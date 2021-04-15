package pwapi

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/internal/randstring"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
)

func (svc *service) UserGetSession(ctx context.Context, _ *UserGetSession_Input) (*UserGetSession_Output, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}
	svc.logger.Debug("token", zap.Any("token", token))

	output := &UserGetSession_Output{}
	output.Claims = pwsso.ClaimsFromToken(token)

	// try loading it from database
	output.User, err = svc.loadOAuthUser(output.Claims.ActionToken.Sub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.ErrGetOAuthUser.Wrap(err)
	}

	// new user
	if errors.Is(err, gorm.ErrRecordNotFound) {
		output.IsNewUser = true
		if _, err = svc.newUserFromClaims(output.Claims); err != nil {
			return nil, errcode.ErrNewUserFromClaims.Wrap(err)
		}
		if output.User, err = svc.loadOAuthUser(output.Claims.ActionToken.Sub); err != nil {
			return nil, errcode.ErrGetOAuthUser.Wrap(err)
		}
	}

	if output.User.Username != output.Claims.PreferredUsername {
		if err := svc.db.Model(output.User).Updates(pwdb.User{Username: output.Claims.PreferredUsername}).Error; err != nil {
			return nil, pwdb.GormToErrcode(err)
		}
		// FIXME: also update the solo organization
	}

	// FIXME: output.Notifications = COUNT
	output.Notifications = int32(rand.Intn(10))

	output.Seasons, err = svc.loadUserSeasons(ctx)
	if err != nil {
		return nil, errcode.ErrLoadUserSeasons.Wrap(err)
	}

	/* FIXME: having a login activity would be nice, but we need to find a solution to detect a login vs a simple user get session refresh
	if !output.IsNewUser {
		activity := pwdb.Activity{
			Kind:     pwdb.Activity_UserLogin,
			AuthorID: output.User.ID,
			UserID:   output.User.ID,
		}
		if err := svc.db.Create(&activity).Error; err != nil {
			return nil, pwdb.GormToErrcode(err)
		}
	}
	*/

	return output, nil
}

func (svc *service) loadUserSeasons(ctx context.Context) ([]*UserGetSession_Output_SeasonAndTeam, error) {
	var (
		seasons     []*pwdb.Season
		memberships []*pwdb.TeamMember
	)

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated
	}

	// get season organizations for user
	err = svc.db.
		Preload("Team").
		Preload("Team.Organization").
		Where(pwdb.TeamMember{UserID: userID}).
		Find(&memberships).
		Error
	if err != nil {
		return nil, errcode.ErrGetUserOrganizations.Wrap(err)
	}

	req := svc.db
	switch {
	case isAdminContext(ctx): // because it's the highest level
		// noop
	case isTesterContext(ctx):
		req = req.
			Where(pwdb.Season{Visibility: pwdb.Season_Public}).
			Or(pwdb.Season{IsTesting: true})
	default: // "normal" user
		req = req.
			Where(pwdb.Season{Visibility: pwdb.Season_Public})
	}

	// get all available seasons
	err = req.
		Find(&seasons).
		Error
	if err != nil {
		return nil, errcode.ErrGetSeasons.Wrap(err)
	}

	output := []*UserGetSession_Output_SeasonAndTeam{}
	for _, season := range seasons {
		item := &UserGetSession_Output_SeasonAndTeam{
			Season: season,
		}

		for _, membership := range memberships {
			if membership.Team.SeasonID == season.ID {
				item.Team = membership.Team
				break
			}
		}

		output = append(output, item)
	}

	return output, nil
}

func (svc *service) loadOAuthUser(subject string) (*pwdb.User, error) {
	var user pwdb.User
	err := svc.db.
		Preload("ActiveTeamMember").
		Preload("ActiveTeamMember.Team").
		Preload("ActiveTeamMember.Team.Season").
		Preload("ActiveTeamMember.Team.Organization").
		Where(pwdb.User{OAuthSubject: subject}).
		First(&user).
		Error
	if err != nil {
		return nil, errcode.ErrGetUserBySubject.Wrap(err)
	}

	return &user, nil
}

func (svc *service) newUserFromClaims(claims *pwsso.Claims) (*pwdb.User, error) {
	if !claims.EmailVerified {
		return nil, errcode.ErrEmailAddressNotVerified
	}

	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(claims.Email)))

	var season pwdb.Season
	if err := svc.db.Where(pwdb.Season{IsGlobal: true}).First(&season).Error; err != nil {
		return nil, errcode.ErrGetDefaultSeason.Wrap(err)
	}

	username := claims.PreferredUsername
	if claims.PreferredUsername == claims.Email {
		username = strings.Split(claims.PreferredUsername, "@")[0]
	}
	username = strings.TrimSpace(username)
	if username == "" {
		username = randstring.RandString(8)
	}

	user := pwdb.User{
		Username:     username,
		Email:        claims.Email,
		GravatarURL:  gravatarURL,
		OAuthSubject: claims.ActionToken.Sub,
		// WebsiteURL
		// Locale

		TeamMemberships:         []*pwdb.TeamMember{},
		OrganizationMemberships: []*pwdb.OrganizationMember{},
		DeletionStatus:          pwdb.DeletionStatus_Active,
	}
	organization := pwdb.Organization{
		Name:           username,
		GravatarURL:    gravatarURL,
		DeletionStatus: pwdb.DeletionStatus_Active,
		GlobalSeason:   true,
		// Locale
	}
	organizationMember := pwdb.OrganizationMember{
		// User: &user,
		Organization: &organization,
		Role:         pwdb.OrganizationMember_Owner,
		Slug:         slug.Make(user.Username),
	}
	team := pwdb.Team{
		Season:         &season,
		IsGlobal:       true,
		Organization:   &organization,
		DeletionStatus: pwdb.DeletionStatus_Active,
	}
	seasonMember := pwdb.TeamMember{
		User: &user,
		Team: &team,
		Role: pwdb.TeamMember_Owner,
		Slug: slug.Make(user.Username),
	}
	user.OrganizationMemberships = []*pwdb.OrganizationMember{&organizationMember}

	err := svc.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&seasonMember).Error; err != nil {
			return err
		}
		activity := pwdb.Activity{
			Kind:           pwdb.Activity_UserRegister,
			AuthorID:       user.ID, // FIXME: author should be based on ctx
			UserID:         user.ID,
			TeamID:         team.ID,
			OrganizationID: organization.ID,
			TeamMemberID:   seasonMember.ID,
			SeasonID:       season.ID,
		}
		if err := tx.Create(&activity).Error; err != nil {
			return err
		}

		// set active season
		err := tx.Model(&user).Updates(pwdb.User{
			ActiveTeamMemberID: seasonMember.ID,
			ActiveSeasonID:     season.ID,
		}).
			Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errcode.ErrCommitUserTransaction.Wrap(err)
	}

	return &user, nil
}
