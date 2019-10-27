package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) TeamGet(ctx context.Context, in *TeamGetInput) (*TeamGetOutput, error) {
	{ // validation
		if in.TeamID == 0 {
			return nil, ErrMissingArgument
		}
	}

	var item pwdb.Team
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.Team{ID: in.TeamID}).
		First(&item).
		Error

	switch {
	case err != nil && pwdb.IsRecordNotFoundError(err):
		return nil, ErrInvalidArgument // FIXME: wrap original error
	case err != nil:
		return nil, fmt.Errorf("fetch team from db: %w", err)
	}

	ret := TeamGetOutput{
		Item: &item,
	}

	return &ret, nil
}
