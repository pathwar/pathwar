package server

import (
	"context"

	"pathwar.pw/entity"
)

func (s *svc) Tournaments(ctx context.Context, _ *Void) (*entity.TournamentList, error) {
	var tournaments entity.TournamentList
	if err := s.db.Set("gorm:auto_preload", true).Find(&tournaments.Items).Error; err != nil {
		return nil, err
	}

	return &tournaments, nil
}
