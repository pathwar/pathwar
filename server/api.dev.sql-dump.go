package server

import (
	"context"

	"pathwar.pw/entity"
	"pathwar.pw/sql"
)

func (s *svc) Dump(ctx context.Context, _ *Void) (*entity.Dump, error) {
	return sql.DoDump(s.db)
}
