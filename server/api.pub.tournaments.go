package server

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/entity"
)

func (s *svc) tournaments(ctx context.Context, _ *Void) ([]*UserSessionOutput_Tournament, error) {
	var tournaments []*entity.Tournament
	var memberships []*entity.TournamentMember

	userID, err := subjectFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.db.
		Where(entity.Tournament{Visibility: entity.Tournament_Public}). // FIXME: admin can see everything
		Find(&tournaments).
		Error; err != nil {
		return nil, err
	}

	// FIXME: should be doable in a unique request with LEFT joining
	if err := s.db.
		Preload("TournamentTeam").
		Preload("TournamentTeam.Team").
		Where(entity.TournamentMember{UserID: userID}).
		Find(&memberships).
		Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	output := []*UserSessionOutput_Tournament{}

	for _, tournament := range tournaments {
		item := &UserSessionOutput_Tournament{
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
