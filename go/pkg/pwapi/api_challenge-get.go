package pwapi

import (
	"context"

	"pathwar.land/go/v2/pkg/errcode"
	"pathwar.land/go/v2/pkg/pwdb"
)

func (svc *service) ChallengeGet(ctx context.Context, in *ChallengeGet_Input) (*ChallengeGet_Output, error) {
	// validation
	if in == nil || in.ChallengeID == 0 {
		return nil, errcode.ErrMissingInput
	}

	var item pwdb.Challenge
	err := svc.db.
		Preload("Flavors").
		Where(pwdb.Challenge{ID: in.ChallengeID}).
		First(&item).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	ret := ChallengeGet_Output{
		Item: &item,
	}
	return &ret, nil
}
