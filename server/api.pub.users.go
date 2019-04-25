package server

import (
	"context"

	"pathwar.pw/entity"
)

func (s *svc) Users(ctx context.Context, _ *Void) (*entity.UserList, error) {
	var users entity.UserList
	if err := s.db.Find(&users.Items).Error; err != nil {
		return nil, err
	}

	return &users, nil
}
