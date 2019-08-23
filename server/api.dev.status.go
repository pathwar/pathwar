package server

import (
	"context"

	"pathwar.land/entity"
)

func (s *svc) Status(ctx context.Context, _ *Void) (*entity.Status, error) {
	return &entity.Status{
		EverythingIsOK: true,
	}, nil
}
