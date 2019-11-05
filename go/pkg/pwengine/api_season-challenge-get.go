package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeGet(ctx context.Context, in *SeasonChallengeGetInput) (*SeasonChallengeGetOutput, error) {
	{ // validation
		if in.SeasonChallengeID == 0 {
			return nil, ErrMissingArgument
		}
	}

	var item pwdb.SeasonChallenge
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.SeasonChallenge{ID: in.SeasonChallengeID}).
		First(&item).
		Error

	switch {
	case err != nil && pwdb.IsRecordNotFoundError(err):
		return nil, ErrInvalidArgument // FIXME: wrap original error
	case err != nil:
		return nil, fmt.Errorf("query season challenge: %w", err)
	}

	ret := SeasonChallengeGetOutput{
		Item: &item,
	}

	return &ret, nil
}
