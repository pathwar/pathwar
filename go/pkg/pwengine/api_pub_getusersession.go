package pwengine

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func (e *engine) GetUserSession(ctx context.Context, _ *Void) (*UserSessionOutput, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token from context")
	}
	zap.L().Debug("token", zap.Any("token", token))

	output := &UserSessionOutput{}
	output.Claims = pwsso.ClaimsFromToken(token)

	// try loading it from database
	output.User, err = e.loadOAuthUser(output.Claims.ActionToken.Sub)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		// internal error
		return nil, err
	}

	// new user
	if gorm.IsRecordNotFoundError(err) {
		output.IsNewUser = true
		if _, err = e.newUserFromClaims(output.Claims); err != nil {
			return nil, err
		}
		if output.User, err = e.loadOAuthUser(output.Claims.ActionToken.Sub); err != nil {
			return nil, err
		}
	}

	if output.User.Username != output.Claims.PreferredUsername {
		return nil, fmt.Errorf("username differs from JWT token and database")
	}

	// FIXME: output.Notifications = COUNT
	output.Notifications = 42

	output.Tournaments, err = e.tournaments(ctx)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (e *engine) tournaments(ctx context.Context) ([]*UserSessionOutput_TournamentAndTeam, error) {
	var (
		tournaments []*pwdb.Tournament
		memberships []*pwdb.TournamentMember
	)

	userID, err := subjectFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := e.db.
		Where(pwdb.Tournament{Visibility: pwdb.Tournament_Public}). // FIXME: admin can see everything
		Find(&tournaments).
		Error; err != nil {
		return nil, err
	}

	// FIXME: should be doable in a unique request with LEFT joining
	if err := e.db.
		Preload("TournamentTeam").
		Preload("TournamentTeam.Team").
		Where(pwdb.TournamentMember{UserID: userID}).
		Find(&memberships).
		Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	output := []*UserSessionOutput_TournamentAndTeam{}

	for _, tournament := range tournaments {
		item := &UserSessionOutput_TournamentAndTeam{
			Tournament: tournament,
		}

		for _, membership := range memberships {
			if membership.TournamentTeam.TournamentID == tournament.ID {
				item.Team = membership.TournamentTeam
				break
			}
		}

		output = append(output, item)
	}

	return output, nil
}

func (e *engine) loadOAuthUser(subject string) (*pwdb.User, error) {
	var user pwdb.User
	if err := e.db.
		Preload("ActiveTournamentMember").
		Preload("ActiveTournamentMember.TournamentTeam").
		Preload("ActiveTournamentMember.TournamentTeam.Tournament").
		Preload("ActiveTournamentMember.TournamentTeam.Team").
		Where(pwdb.User{ID: subject}).
		First(&user).
		Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (e *engine) newUserFromClaims(claims *pwsso.Claims) (*pwdb.User, error) {
	if claims.EmailVerified == false {
		return nil, fmt.Errorf("you need to verify your email address")
	}

	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(claims.Email)))

	var tournament pwdb.Tournament
	if err := e.db.Where(pwdb.Tournament{IsDefault: true}).First(&tournament).Error; err != nil {
		return nil, err
	}

	user := pwdb.User{
		ID:          claims.ActionToken.Sub,
		Username:    claims.PreferredUsername,
		Email:       claims.Email,
		GravatarURL: gravatarURL,
		// WebsiteURL
		// Locale

		TournamentMemberships: []*pwdb.TournamentMember{},
		Memberships:           []*pwdb.TeamMember{},
	}
	team := pwdb.Team{
		Name:        claims.PreferredUsername,
		GravatarURL: gravatarURL,
		// Locale
	}
	teamMember := pwdb.TeamMember{
		//User: &user,
		Team: &team,
		Role: pwdb.TeamMember_Owner,
	}
	tournamentTeam := pwdb.TournamentTeam{
		Tournament: &tournament,
		IsDefault:  true,
		Team:       &team,
	}
	tournamentMember := pwdb.TournamentMember{
		User:           &user,
		TournamentTeam: &tournamentTeam,
		Role:           pwdb.TournamentMember_Owner,
	}
	user.Memberships = []*pwdb.TeamMember{&teamMember}

	tx := e.db.Begin()
	tx.Create(&user)
	tx.Create(&tournamentMember)

	// FIXME: create a "welcome" notification

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// set active tournament
	if err := e.db.
		Model(&user).
		Updates(pwdb.User{
			ActiveTournamentMemberID: tournamentMember.ID,
			ActiveTournamentID:       tournament.ID,
		}).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}
