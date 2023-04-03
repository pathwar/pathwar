package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

// AdminUpdateSeasonMetadata Actually, this update only the subscription field of the season object but it could be extended to update other fields in the future
func (svc *service) AdminUpdateSeasonMetadata(ctx context.Context, in *AdminUpdateSeasonMetadata_Input) (*AdminUpdateSeasonMetadata_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.Season == nil {
		return nil, errcode.ErrMissingInput
	}

	err := svc.db.Model(in.Season).Update("subscription", in.Season.Subscription).Error
	if err != nil {
		return nil, err
	}

	return &AdminUpdateSeasonMetadata_Output{Season: in.Season}, nil
}
