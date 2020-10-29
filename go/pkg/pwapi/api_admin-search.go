package pwapi

import (
	"context"
	"fmt"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminSearch(ctx context.Context, in *AdminSearch_Input) (*AdminSearch_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.Search == "" {
		return nil, errcode.ErrMissingInput
	}

	dbLikeSearchStr := fmt.Sprint("%", in.Search, "%")

	query := svc.db.Where("id LIKE ? OR slug LIKE ?", dbLikeSearchStr, dbLikeSearchStr)

	out := AdminSearch_Output{}
	err := query.
		Or("name LIKE ?", dbLikeSearchStr).
		Find(&out.Challenges).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Or("category LIKE ?", dbLikeSearchStr).
		Find(&out.ChallengeFlavors).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.SeasonChallenges).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.ChallengeInstances).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Or("name LIKE ? OR hostname LIKE ?", dbLikeSearchStr, dbLikeSearchStr).
		Find(&out.Agents).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.OrganizationMembers).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.TeamMembers).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Or("username LIKE ? OR email LIKE ?", dbLikeSearchStr, dbLikeSearchStr).
		Find(&out.Users).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Or("name LIKE ?", dbLikeSearchStr).
		Find(&out.Organizations).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Or("name LIKE ?", dbLikeSearchStr).
		Find(&out.Seasons).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.Teams).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.WhoswhoAttempts).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.ChallengeValidations).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.ChallengeSubscriptions).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.InventoryItems).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = query.
		Find(&out.Achievements).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Where("id LIKE ?", dbLikeSearchStr).
		Find(&out.Notifications).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Where("id LIKE ?", dbLikeSearchStr).
		Find(&out.Coupons).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Where("id LIKE ?", dbLikeSearchStr).
		Find(&out.CouponValidations).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Where("id LIKE ?", dbLikeSearchStr).
		Find(&out.Activities).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	return &out, nil
}
