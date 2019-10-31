package pwengine

import (
	"context"
	"crypto/md5"
	"fmt"

	"go.uber.org/zap"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func (e *engine) UserGetSession(ctx context.Context, _ *UserGetSessionInput) (*UserGetSessionOutput, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get token from context: %w", err)
	}
	zap.L().Debug("token", zap.Any("token", token))

	output := &UserGetSessionOutput{}
	output.Claims = pwsso.ClaimsFromToken(token)

	// try loading it from database
	output.User, err = e.loadOAuthUser(output.Claims.ActionToken.Sub)
	if err != nil && !pwdb.IsRecordNotFoundError(err) {
		// internal error
		return nil, fmt.Errorf("load oauth user: %w", err)
	}

	// new user
	if pwdb.IsRecordNotFoundError(err) {
		output.IsNewUser = true
		if _, err = e.newUserFromClaims(output.Claims); err != nil {
			return nil, fmt.Errorf("new user from claims: %w", err)
		}
		if output.User, err = e.loadOAuthUser(output.Claims.ActionToken.Sub); err != nil {
			return nil, fmt.Errorf("load oauth user: %w", err)
		}
	}

	if output.User.Username != output.Claims.PreferredUsername {
		return nil, fmt.Errorf("username differs from JWT token and database")
	}

	// FIXME: output.Notifications = COUNT
	output.Notifications = 42

	output.Seasons, err = e.seasons(ctx)
	if err != nil {
		return nil, fmt.Errorf("get seasons: %w", err)
	}

	return output, nil
}

func (e *engine) seasons(ctx context.Context) ([]*UserGetSessionOutput_SeasonAndTeam, error) {
	var (
		seasons     []*pwdb.Season
		memberships []*pwdb.TeamMember
	)

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	// get season organizations for user
	err = e.db.
		Preload("Team").
		Preload("Team.Organization").
		Where(pwdb.TeamMember{UserID: userID}).
		Find(&memberships).
		Error
	if err != nil && !pwdb.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("fetch season organizations: %w", err)
	}

	// get all available seasons
	err = e.db.
		Where(pwdb.Season{Visibility: pwdb.Season_Public}).
		// FIXME: admins can see everything
		Find(&seasons).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch seasons: %w", err)
	}

	output := []*UserGetSessionOutput_SeasonAndTeam{}
	for _, season := range seasons {
		item := &UserGetSessionOutput_SeasonAndTeam{
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

func (e *engine) loadOAuthUser(subject string) (*pwdb.User, error) {
	var user pwdb.User
	err := e.db.
		Preload("ActiveTeamMember").
		Preload("ActiveTeamMember.Team").
		Preload("ActiveTeamMember.Team.Season").
		Preload("ActiveTeamMember.Team.Organization").
		Where(pwdb.User{OAuthSubject: subject}).
		First(&user).
		Error

	if err != nil {
		return nil, fmt.Errorf("fetch user from subject %q: %w", subject, err)
	}

	return &user, nil
}

func (e *engine) newUserFromClaims(claims *pwsso.Claims) (*pwdb.User, error) {
	if claims.EmailVerified == false {
		return nil, fmt.Errorf("you need to verify your email address")
	}

	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(claims.Email)))

	var season pwdb.Season
	if err := e.db.Where(pwdb.Season{IsDefault: true}).First(&season).Error; err != nil {
		return nil, fmt.Errorf("get default season: %w", err)
	}

	user := pwdb.User{
		Username:     claims.PreferredUsername,
		Email:        claims.Email,
		GravatarURL:  gravatarURL,
		OAuthSubject: claims.ActionToken.Sub,
		// WebsiteURL
		// Locale

		TeamMemberships: []*pwdb.TeamMember{},
		Memberships:     []*pwdb.OrganizationMember{},
	}
	organization := pwdb.Organization{
		Name:        claims.PreferredUsername,
		GravatarURL: gravatarURL,
		// Locale
	}
	organizationMember := pwdb.OrganizationMember{
		//User: &user,
		Organization: &organization,
		Role:         pwdb.OrganizationMember_Owner,
	}
	seasonOrganization := pwdb.Team{
		Season:       &season,
		IsDefault:    true,
		Organization: &organization,
	}
	seasonMember := pwdb.TeamMember{
		User: &user,
		Team: &seasonOrganization,
		Role: pwdb.TeamMember_Owner,
	}
	user.Memberships = []*pwdb.OrganizationMember{&organizationMember}

	tx := e.db.Begin()
	tx.Create(&user)
	tx.Create(&seasonMember)

	// FIXME: create a "welcome" notification

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// set active season
	err := e.db.
		Model(&user).
		Updates(pwdb.User{
			ActiveTeamMemberID: seasonMember.ID,
			ActiveSeasonID:     season.ID,
		}).
		Error
	if err != nil {
		return nil, fmt.Errorf("update active season: %w", err)
	}

	return &user, nil
}
