package pwapi

import (
	"context"
	"fmt"

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

	fmt.Println(in)
	var activities []*pwdb.Activity
	req := svc.db.
		Preload("Author").
		Preload("Team").
		Preload("User").
		Preload("Agent").
		Preload("Organization").
		Preload("Season").
		Preload("Challenge").
		Preload("ChallengeFlavor").
		Preload("ChallengeInstance").
		Preload("Coupon").
		Preload("SeasonChallenge").
		Preload("TeamMember").
		Preload("ChallengeSubscription")
	if in.Limit > 0 {
		fmt.Println("TEST, in.Limit")
		req = req.Limit(in.Limit)
	}
	if in.Since != nil {
		return nil, errcode.ErrNotImplemented
	}
	switch in.FilteringPreset {
	case "default", "":
	// noop
	case "registers":
		req = req.Where(&pwdb.Activity{Kind: pwdb.Activity_UserRegister})
	default:
		return nil, errcode.TODO
	}

	if err := req.Find(&activities).Error; err != nil {
		return nil, errcode.ErrListActivities.Wrap(err)
	}

	out := AdminListActivities_Output{Activities: activities}
	return &out, nil
}
