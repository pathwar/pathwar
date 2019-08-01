package server

import (
	"context"
	"time"

	"pathwar.land/entity"
	"pathwar.land/version"
)

func (s *svc) Info(ctx context.Context, _ *Void) (*entity.Info, error) {
	return &entity.Info{
		Version: version.Version,
		Commit:  version.Commit,
		BuiltAt: version.Date,
		BuiltBy: version.BuiltBy,
		Uptime:  int32(time.Now().Sub(s.startedAt).Seconds()),
	}, nil
}
