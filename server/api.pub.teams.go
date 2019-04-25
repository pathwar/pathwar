package server

import (
	"context"

	"pathwar.pw/entity"
)

func (s *svc) Teams(ctx context.Context, _ *Void) (*entity.TeamList, error) {
	var teams entity.TeamList
	if err := s.db.Set("gorm:auto_preload", true).Find(&teams.Items).Error; err != nil {
		return nil, err
	}

	// FIXME: filter-out
	return &teams, nil
}
