package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeAdd(ctx context.Context, in *AdminChallengeAdd_Input) (*AdminChallengeAdd_Output, error) {
	// FIXME: check if client is admin

	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var challenge pwdb.Challenge
	challenge.Name = in.Name
	challenge.Description = in.Description
	challenge.Author = in.Author
	challenge.Locale = in.Locale
	challenge.IsDraft = in.IsDraft
	challenge.PreviewUrl = in.PreviewUrl
	challenge.Homepage = in.Homepage

	err := svc.db.Create(&challenge).Error
	if err != nil {
		return nil, errcode.ErrChallengeAdd.Wrap(err)
	}

	out := AdminChallengeAdd_Output{
		Challenge: &challenge,
	}
	return &out, nil
}
