package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListCoupons(ctx context.Context, in *AdminListCoupons_Input) (*AdminListCoupons_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var coupons []*pwdb.Coupon
	err := svc.db.
		Preload("Season").
		Preload("Validations").
		Preload("Validations.Author").
		Preload("Validations.Team").
		Preload("Validations.Team.Organization").
		Find(&coupons).Error
	if err != nil {
		return nil, errcode.ErrListCoupons.Wrap(err)
	}

	out := AdminListCoupons_Output{Coupons: coupons}
	return &out, nil
}
