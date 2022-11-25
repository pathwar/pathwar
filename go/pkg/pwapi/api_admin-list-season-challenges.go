package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminListSeasonChallenges(ctx context.Context, in *AdminListSeasonChallenges_Input) (*AdminListSeasonChallenges_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	return &AdminListSeasonChallenges_Output{}, nil
}
