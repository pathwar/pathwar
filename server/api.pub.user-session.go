package server

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"pathwar.land/client"
	"pathwar.land/entity"
)

func (s *svc) UserSession(ctx context.Context, _ *Void) (*entity.UserSession, error) {
	token, err := userTokenFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token from context")
	}
	zap.L().Debug("token", zap.Any("token", token))

	sess, err := client.UserSessionFromToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session from token")
	}
	return &sess, nil
}
