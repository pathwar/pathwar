package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/internal/randstring"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminAddCoupon(ctx context.Context, in *AdminAddCoupon_Input) (*AdminAddCoupon_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	in.ApplyDefaults()
	if in == nil || in.Hash == "" || in.Value == 0 || in.SeasonID == "" || in.MaxValidationCount == 0 {
		return nil, errcode.ErrMissingInput
	}

	if in.Hash == "RANDOM" {
		in.Hash = randstring.RandString(16)
	}

	seasonID, err := pwdb.GetIDBySlugAndKind(svc.db, in.SeasonID, "season")
	if err != nil {
		return nil, err
	}

	coupon := pwdb.Coupon{
		Hash:               in.Hash,
		Value:              in.Value,
		SeasonID:           seasonID,
		MaxValidationCount: in.MaxValidationCount,
	}
	err = svc.db.Create(&coupon).Error
	if err != nil {
		return nil, errcode.ErrAddCoupon.Wrap(err)
	}

	out := AdminAddCoupon_Output{
		Coupon: &coupon,
	}
	return &out, nil
}

func (in *AdminAddCoupon_Input) ApplyDefaults() {
	if in.Hash == "" {
		in.Hash = "RANDOM"
	}
	if in.MaxValidationCount == 0 {
		in.MaxValidationCount = 1
	}
	if in.SeasonID == "" {
		in.SeasonID = "solo-mode"
	}
	if in.Value == 0 {
		in.Value = 1
	}
}
