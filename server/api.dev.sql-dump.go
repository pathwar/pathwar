package server

import (
	"context"

	"pathwar.land/entity"
	"pathwar.land/sql"
)

func (s *svc) Dump(ctx context.Context, _ *Void) (*entity.Dump, error) {
	return sql.DoDump(s.db)
}
