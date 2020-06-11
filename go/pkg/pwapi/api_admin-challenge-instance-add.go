package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminChallengeInstanceAdd(ctx context.Context, in *AdminChallengeInstanceAdd_Input) (*AdminChallengeInstanceAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var agent pwdb.Agent
	err := svc.db.
		Where(pwdb.Agent{ID: in.ChallengeInstance.AgentID}).
		First(&agent).
		Error
	if err != nil {
		return nil, errcode.ErrChallengeInstanceAdd.Wrap(err)
	}

	var flavor pwdb.ChallengeFlavor
	err = svc.db.
		Where(pwdb.ChallengeFlavor{ID: in.ChallengeInstance.FlavorID}).
		First(&flavor).
		Error
	if err != nil {
		return nil, errcode.ErrChallengeInstanceAdd.Wrap(err)
	}

	challengeInstance := in.ChallengeInstance

	err = svc.db.Create(&challengeInstance).Error
	if err != nil {
		return nil, errcode.ErrChallengeInstanceAdd.Wrap(err)
	}

	out := AdminChallengeInstanceAdd_Output{
		ChallengeInstance: challengeInstance,
	}
	return &out, nil
}
