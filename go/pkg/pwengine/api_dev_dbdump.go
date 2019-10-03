package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) DBDump(context.Context, *Void) (*pwdb.Dump, error) {
	return pwdb.GetDump(e.db)
}
