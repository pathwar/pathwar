package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) ToolDBDump(context.Context, *Void) (*pwdb.Dump, error) {
	return pwdb.GetDump(svc.db)
}
