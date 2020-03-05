package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeAdd(ctx context.Context, in *AdminChallengeAdd_Input) (*AdminChallengeAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var challenge pwdb.Challenge
	challenge.Name = in.Challenge.Name
	challenge.Description = in.Challenge.Description
	challenge.Author = in.Challenge.Author
	challenge.Locale = in.Challenge.Locale
	challenge.IsDraft = in.Challenge.IsDraft
	challenge.PreviewUrl = in.Challenge.PreviewUrl
	challenge.Homepage = in.Challenge.Homepage

	err := svc.db.Create(&challenge).Error
	if err != nil {
		return nil, errcode.ErrChallengeAdd.Wrap(err)
	}

	out := AdminChallengeAdd_Output{
		Challenge: &challenge,
	}
	return &out, nil
}
