package pwengine

import (
	context "context"

	pwdb "pathwar.land/go/pkg/pwdb"
)

func (c *client) GenerateFakeData(context.Context, *Void) (*Void, error) {
	return &Void{}, pwdb.GenerateFakeData(c.db, c.logger.Named("generate-fake-data"))
}
