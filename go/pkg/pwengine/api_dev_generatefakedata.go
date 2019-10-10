package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) GenerateFakeData(context.Context, *Void) (*Void, error) {
	return &Void{}, pwdb.GenerateFakeData(e.db, e.logger.Named("generate-fake-data"))
}
