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

	var user entity.User
	err = s.db.Where(entity.User{OauthSubject: output.Claims.ActionToken.Sub}).First(&user).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return output, err
	}

	if gorm.IsRecordNotFoundError(err) {
		output.IsNewUser = true
		output.User, err = s.newUser(output.Claims)
		if err != nil {
			return output, err
		}
		// FIXME: reload the user with all the preloads? or recursive call
	} else {
		output.User = &user
		if output.User.Username != output.Claims.PreferredUsername {
			return output, fmt.Errorf("username differs from JWT token and database")
		}
	}

	// FIXME: ActiveTournamentTeam: ...
	// FIXME: ActiveTeam: ...
	// FIXME: ActiveTournament: ...

	// FIXME: Notifications: 0,

	return output, nil
}

func (s *svc) newUser(claims *client.Claims) (*entity.User, error) {
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
	tx.Create(&user)

	// FIXME: create a "welcome" notification

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// set active tournament
	if err := s.db.Model(&user).Updates(entity.User{ActiveTournamentMemberID: tournamentMember.Metadata.ID}).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
