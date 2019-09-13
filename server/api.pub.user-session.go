package server

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"pathwar.land/client"
	"pathwar.land/entity"
)

func (s *svc) UserSession(ctx context.Context, _ *Void) (*UserSessionOutput, error) {
	token, err := userTokenFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token from context")
	}
	zap.L().Debug("token", zap.Any("token", token))

	output := &UserSessionOutput{}
	output.Claims = client.ClaimsFromToken(token)

	// try loading it from database
	output.User, err = s.loadOauthUser(output.Claims.ActionToken.Sub)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		// internal error
		return nil, err
	}

	// new user
	if gorm.IsRecordNotFoundError(err) {
		output.IsNewUser = true
		if _, err = s.newUserFromClaims(output.Claims); err != nil {
			return nil, err
		}
		if output.User, err = s.loadOauthUser(output.Claims.ActionToken.Sub); err != nil {
			return nil, err
		}
	}

	if output.User.Username != output.Claims.PreferredUsername {
		return nil, fmt.Errorf("username differs from JWT token and database")
	}

	// FIXME: output.Notifications = COUNT
	output.Notifications = 42

	return output, nil
}

func (s *svc) loadOauthUser(subject string) (*entity.User, error) {
	var user entity.User
	if err := s.db.
		Preload("ActiveTournamentMember").
		Preload("ActiveTournamentMember.TournamentTeam").
		Preload("ActiveTournamentMember.TournamentTeam.Tournament").
		Preload("ActiveTournamentMember.TournamentTeam.Team").
		Where(entity.User{OauthSubject: subject}).
		First(&user).
		Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *svc) newUserFromClaims(claims *client.Claims) (*entity.User, error) {
	if claims.EmailVerified == false {
		return nil, fmt.Errorf("you need to verify your email address")
	}

	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(claims.Email)), 200) // FIXME: remove size and handle it on the fly

	var tournament entity.Tournament
	if err := s.db.Where(entity.Tournament{IsDefault: true}).First(&tournament).Error; err != nil {
		return nil, err
	}

	team := entity.Team{
		Name:        claims.PreferredUsername,
		GravatarURL: gravatarURL,
		// Locale
	}
	user := entity.User{
		OauthSubject: claims.ActionToken.Sub,
		Username:     claims.PreferredUsername,
		Email:        claims.Email,
		GravatarURL:  gravatarURL,
		// WebsiteURL
		// Locale

		TournamentMemberships: []*entity.TournamentMember{},
		Memberships:           []*entity.TeamMember{},
	}
	teamMember := entity.TeamMember{
		//User: &user,
		Team: &team,
		Role: entity.TeamMember_Owner,
	}
	tournamentTeam := entity.TournamentTeam{
		Tournament: &tournament,
		IsDefault:  true,
		Team:       &team,
	}
	tournamentMember := entity.TournamentMember{
		User:           &user,
		TournamentTeam: &tournamentTeam,
		Role:           entity.TournamentMember_Owner,
	}
	user.Memberships = []*entity.TeamMember{&teamMember}

	tx := s.db.Begin()
	tx.Create(&tournamentMember)

	// FIXME: create a "welcome" notification

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// set active tournament
	if err := s.db.
		Model(&user).
		Updates(entity.User{ActiveTournamentMemberID: tournamentMember.Metadata.ID}).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}
