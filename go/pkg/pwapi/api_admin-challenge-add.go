package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeAdd(ctx context.Context, in *AdminChallengeAdd_Input) (*AdminChallengeAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	in.ApplyDefaults()
	if in == nil || in.Challenge.Name == "" {
		return nil, errcode.ErrMissingInput
	}

	challenge := pwdb.Challenge{
		Slug:        in.Challenge.Slug,
		Name:        in.Challenge.Name,
		Description: in.Challenge.Description,
		Author:      in.Challenge.Author,
		Locale:      in.Challenge.Locale,
		IsDraft:     in.Challenge.IsDraft,
		PreviewURL:  in.Challenge.PreviewURL,
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

func (in *AdminChallengeAdd_Input) ApplyDefaults() {
	if in.Challenge.Locale == "" {
		in.Challenge.Locale = "en_US"
	}
	if in.Challenge.Author == "" {
		in.Challenge.Author = "Pathwar Staff"
	}
	if in.Challenge.Slug != "" && in.Challenge.Name == "" {
		in.Challenge.Name = in.Challenge.Slug
	}
}
