package server

import (
	"context"

	"pathwar.land/entity"
)

func (s *svc) Tournaments(ctx context.Context, _ *Void) (*TournamentsOutput, error) {
	var tournaments entity.TournamentList
	if err := s.db.Set("gorm:auto_preload", true).Find(&tournaments.Items).Error; err != nil {
		return nil, err
	}

	// FIXME: fetch user teams for each tournament

	output := &TournamentsOutput{
		Items: []*TournamentsOutput_Tournament{},
	}

	for _, tournament := range tournaments.Items {
		output.Items = append(
			output.Items,
			&TournamentsOutput_Tournament{
				Tournament: tournament,
				Team:       nil,
			},
		)
	}

	return output, nil
}
