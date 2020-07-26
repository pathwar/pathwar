package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListActivities(ctx context.Context, in *AdminListActivities_Input) (*AdminListActivities_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var activities []*pwdb.Activity
	err := svc.db.
		Preload("Author").
		Preload("Team").
		Preload("User").
		Preload("Organization").
		Preload("Season").
		Preload("Challenge").
		Preload("Coupon").
		Preload("SeasonChallenge").
		Preload("TeamMember").
		Preload("ChallengeSubscription").
		Find(&activities).Error
	if err != nil {
		return nil, errcode.ErrListActivities.Wrap(err)
	}

	out := AdminListActivities_Output{Activities: activities}
	return &out, nil
}
