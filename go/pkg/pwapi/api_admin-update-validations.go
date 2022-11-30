package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminUpdateValidations(ctx context.Context, in *AdminUpdateValidations_Input) (*AdminUpdateValidations_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.SeasonChallenge == nil {
		return nil, errcode.ErrMissingInput
	}

	for _, challenge := range in.SeasonChallenge {
		err := svc.db.Model(challenge).Update("NbValidations", challenge.NbValidations).Error
		if err != nil {
			return nil, err
		}
	}

	return &AdminUpdateValidations_Output{}, nil
}
