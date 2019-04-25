package server

import (
	"context"

	"pathwar.pw/entity"
)

func (s *svc) Levels(ctx context.Context, _ *Void) (*entity.LevelList, error) {
	var levels entity.LevelList
	if err := s.db.Set("gorm:auto_preload", true).Find(&levels.Items).Error; err != nil {
		return nil, err
	}

	return &levels, nil
}
