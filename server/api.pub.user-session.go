package server

import (
	"context"

	"github.com/pkg/errors"

	"pathwar.land/entity"
)

func (s *svc) UserSession(ctx context.Context, _ *Void) (*entity.UserSession, error) {
	sess, err := userSessionFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get context session")
	}
	return &sess, nil
}
