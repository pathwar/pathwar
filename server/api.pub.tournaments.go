package server

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/entity"
)

func (s *svc) Tournaments(ctx context.Context, _ *Void) (*TournamentsOutput, error) {
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

	output := &TournamentsOutput{
		Items: []*TournamentsOutput_Tournament{},
	}

	for _, tournament := range tournaments {
		item := &TournamentsOutput_Tournament{
			Tournament: tournament,
		}

		for _, membership := range memberships {
			if membership.TournamentTeam.TournamentID == tournament.ID {
				item.Team = membership.TournamentTeam
				break
			}
		}

		output.Items = append(output.Items, item)
	}

	return output, nil
}
