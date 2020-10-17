package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListAll(ctx context.Context, in *AdminListAll_Input) (*AdminListAll_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	out := AdminListAll_Output{}
	err := svc.db.
		Find(&out.Challenges).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.ChallengeFlavors).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.SeasonChallenges).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.ChallengeInstances).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Agents).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.OrganizationMembers).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.TeamMembers).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.TeamInvites).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Users).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Organizations).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Seasons).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Teams).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.WhoswhoAttempts).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.ChallengeValidations).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.ChallengeSubscriptions).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.InventoryItems).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Notifications).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Coupons).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.CouponValidations).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Achievements).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}
	err = svc.db.
		Find(&out.Activities).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	return &out, nil
}
