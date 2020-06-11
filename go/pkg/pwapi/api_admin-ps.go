package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminPS(ctx context.Context, in *AdminPS_Input) (*AdminPS_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var instances []*pwdb.ChallengeInstance
	err := svc.db.
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Flavor.SeasonChallenges").
		Find(&instances).Error
	if err != nil {
		return nil, errcode.ErrListChallengeInstances.Wrap(err)
	}

	out := AdminPS_Output{Instances: instances}
	return &out, nil
}
