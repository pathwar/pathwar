package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminSeasonStats(ctx context.Context, in *AdminSeasonStats_Input) (*AdminSeasonStats_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	return &AdminSeasonStats_Output{}, nil
}
