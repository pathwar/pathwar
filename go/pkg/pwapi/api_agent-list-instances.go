package pwapi

import (
	"context"

	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (svc *service) AgentListInstances(ctx context.Context, in *AgentListInstances_Input) (*AgentListInstances_Output, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}
	svc.logger.Debug("token", zap.Any("token", token))

	if in == nil || in.AgentName == "" {
		return nil, errcode.ErrMissingInput
	}
	// FIXME: check if client is agent OR admin

	var agent pwdb.Agent
	err = svc.db.
		Where(&pwdb.Agent{
			Status: pwdb.Agent_Active,
			Name:   in.AgentName,
		}).
		First(&agent).
		Error
	if err != nil {
		return nil, errcode.ErrGetAgent.Wrap(err)
	}
	// FIXME: update last seen

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
