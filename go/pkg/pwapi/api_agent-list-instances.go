package pwapi

import (
	"context"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (svc *service) AgentListInstances(ctx context.Context, in *AgentListInstances_Input) (*AgentListInstances_Output, error) {
	if in == nil || in.AgentID == 0 {
		return nil, errcode.ErrMissingInput
	}
	// FIXME: check if client is agent OR admin

	var agent pwdb.Agent
	err := svc.db.
		Where(&pwdb.Agent{Status: pwdb.Agent_Active}).
		First(&agent, in.AgentID).
		Error
	if err != nil {
		return nil, errcode.ErrGetAgent.Wrap(err)
	}
	// FIXME: update last seen

	var instances []*pwdb.ChallengeInstance
	err = svc.db.
		Where(pwdb.ChallengeInstance{AgentID: in.AgentID}). // FIXME: status is active
		Preload("Agent").
		Preload("Flavor").
		Preload("Flavor.SeasonChallenges").
		Preload("Flavor.Challenge").
		Find(&instances).Error
	if err != nil {
		return nil, errcode.ErrListChallengeInstances.Wrap(err)
	}

	out := AgentListInstances_Output{Instances: instances}
	return &out, nil
}
