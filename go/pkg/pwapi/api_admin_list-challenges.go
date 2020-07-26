package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListChallenges(ctx context.Context, in *AdminListChallenges_Input) (*AdminListChallenges_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var challenges []*pwdb.Challenge
	err := svc.db.
		Preload("Flavors").
		Preload("Flavors.Instances").
		Preload("Flavors.SeasonChallenges").
		Preload("Flavors.SeasonChallenges.Season").
		Find(&challenges).Error
	if err != nil {
		return nil, errcode.ErrListChallengeInstances.Wrap(err)
	}

	out := AdminListChallenges_Output{Challenges: challenges}
	return &out, nil
}
