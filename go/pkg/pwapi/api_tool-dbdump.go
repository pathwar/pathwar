package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) ToolDBDump(context.Context, *Void) (*pwdb.Dump, error) {
	return pwdb.GetDump(svc.db)
}
