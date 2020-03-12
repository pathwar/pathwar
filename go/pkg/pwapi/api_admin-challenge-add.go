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

	challenge := pwdb.Challenge{
		Name:        in.Challenge.Name,
		Description: in.Challenge.Description,
		Author:      in.Challenge.Author,
		Locale:      in.Challenge.Locale,
		IsDraft:     in.Challenge.IsDraft,
		PreviewUrl:  in.Challenge.PreviewUrl,
		Homepage:    in.Challenge.Homepage,
	}

	err := svc.db.Create(&challenge).Error
	if err != nil {
		return nil, errcode.ErrChallengeAdd.Wrap(err)
	}

	out := AdminChallengeAdd_Output{
		Challenge: &challenge,
	}
	return &out, nil
}
