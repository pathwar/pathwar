package pwengine

import (
	context "context"

	pwdb "pathwar.land/go/pkg/pwdb"
)

func (c *client) DBDump(context.Context, *Void) (*pwdb.Dump, error) {
	return pwdb.GetDump(c.db)
}
