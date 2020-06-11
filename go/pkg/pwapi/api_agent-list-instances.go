package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AgentListInstances(ctx context.Context, in *AgentListInstances_Input) (*AgentListInstances_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil || in.AgentName == "" {
		return nil, errcode.ErrMissingInput
	}

	var agent pwdb.Agent
	err := svc.db.
		Where(&pwdb.Agent{Name: in.AgentName}).
		First(&agent).
		Error
	if err != nil {
		return nil, errcode.ErrGetAgent.Wrap(err)
	}

	if agent.Status != pwdb.Agent_Active {
		return nil, errcode.ErrInactiveAgent
	}

	// FIXME: update lastSeen and timesSeen

	var instances []*pwdb.ChallengeInstance
	err = svc.db.
		Where(pwdb.ChallengeInstance{AgentID: agent.ID}). // FIXME: status is active
		Preload("Agent").
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Flavor.SeasonChallenges").
		Preload("Flavor.SeasonChallenges.Subscriptions", pwdb.ChallengeSubscription{Status: pwdb.ChallengeSubscription_Active}).
		Preload("Flavor.SeasonChallenges.Subscriptions.Team").
		Preload("Flavor.SeasonChallenges.Subscriptions.Team.Members").
		Find(&instances).Error
	if err != nil {
		return nil, errcode.ErrListChallengeInstances.Wrap(err)
	}

	out := AgentListInstances_Output{Instances: instances}
	return &out, nil
}
