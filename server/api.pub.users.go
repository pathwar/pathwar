package server

import (
	"context"

	"pathwar.land/entity"
)

func (s *svc) Users(ctx context.Context, _ *Void) (*entity.UserList, error) {
	var users entity.UserList
	if err := s.db.Set("gorm:auto_preload", true).Find(&users.Items).Error; err != nil {
		return nil, err
	}

	return &users, nil
}
