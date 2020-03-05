package pwapi

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) CouponValidate(ctx context.Context, in *CouponValidate_Input) (*CouponValidate_Output, error) {
	// validation
	if in == nil || in.Hash == "" || in.TeamID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	// FIXME: create transaction

	// check if user belongs to team
	// FIXME: or is admin
	var team pwdb.Team
	err = svc.db.
		Joins("JOIN team_member ON team_member.team_id = team.id AND team_member.user_id = ?", userID).
		Preload("Members").
		First(&team, in.TeamID).
		Error
	if err != nil {
		return nil, errcode.ErrUserDoesNotBelongToTeam.Wrap(err)
	}

	// fetch coupon
	var coupon pwdb.Coupon
	err = svc.db.
		Where(pwdb.Coupon{Hash: in.Hash, SeasonID: team.SeasonID}).
		First(&coupon).
		Error
	if err != nil {
		return nil, errcode.ErrCouponNotFound.Wrap(err)
	}

	// is already validated by same team
	var validations int32
	err = svc.db.
		Model(&pwdb.CouponValidation{}).
		Where(&pwdb.CouponValidation{CouponID: coupon.ID, TeamID: team.ID}).
		Count(&validations).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	if validations > 0 {
		return nil, errcode.ErrCouponAlreadyValidatedBySameTeam
	}

	// is expired
	err = svc.db.
		Model(&pwdb.CouponValidation{}).
		Where(&pwdb.CouponValidation{CouponID: coupon.ID}).
		Count(&validations).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	if validations >= coupon.MaxValidationCount {
		return nil, errcode.ErrCouponExpired
	}

	// FIXME: validate team
	// FIXME: inacitve user/team

	// create validation
	validation := pwdb.CouponValidation{
		Comment:  "xxx",
		AuthorID: userID,
		TeamID:   team.ID,
		CouponID: coupon.ID,
	}
	err = svc.db.Create(&validation).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	// update team cash
	err = svc.db.Model(&team).UpdateColumn("cash", gorm.Expr("cash + ?", coupon.Value)).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	// load it again with preload
	err = svc.db.
		Preload("Team").
		Preload("Author").
		Preload("Coupon").
		First(&validation, validation.ID).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	ret := CouponValidate_Output{
		CouponValidation: &validation,
	}
	return &ret, nil
}
