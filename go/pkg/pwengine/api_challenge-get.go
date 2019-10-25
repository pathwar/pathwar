package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ChallengeGet(ctx context.Context, in *ChallengeGetInput) (*ChallengeGetOutput, error) {
	{ // validation
		if in.ChallengeID == 0 {
			return nil, ErrMissingArgument
		}
	}

	var item pwdb.Challenge
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.Challenge{ID: in.ChallengeID}).
		First(&item).
		Error

	switch {
	case err != nil && pwdb.IsRecordNotFoundError(err):
		return nil, ErrInvalidArgument // FIXME: wrap original error
	case err != nil:
		return nil, fmt.Errorf("query challenge: %w", err)
	}

	ret := ChallengeGetOutput{
		Item: &item,
	}

	return &ret, nil
}
