package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminUpdateValidations(ctx context.Context, in *AdminUpdateValidations_Input) (*AdminUpdateValidations_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	return nil, nil
}
