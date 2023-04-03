package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminUpdateSeasonMetadata(ctx context.Context, in *AdminUpdateSeasonMetadata_Input) (*AdminUpdateSeasonMetadata_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.Season == nil {
		return nil, errcode.ErrMissingInput
	}

	/*	updates := pwdb.Season{
			IsPublic: in.Season.IsPublic,
		}
		err := svc.db.Model(in.Season).Update(&updates).Error
		if err != nil {
			return nil, err
		}*/

	return &AdminUpdateSeasonMetadata_Output{}, nil
}
