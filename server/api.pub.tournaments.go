package server

import (
	"context"

	"pathwar.pw/entity"
)

func (s *svc) Tournaments(ctx context.Context, _ *Void) (*entity.Tournaments, error) {
	var tournaments entity.Tournaments
	if err := s.db.Find(&tournaments.Items).Error; err != nil {
		return nil, err
	}

	return &tournaments, nil
}
