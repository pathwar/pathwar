package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AgentUpdateState(ctx context.Context, in *AgentUpdateState_Input) (*AgentUpdateState_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	for _, challengeInstance := range in.Instances {
		cpy := challengeInstance
		err := svc.db.Model(&cpy).
			Update(pwdb.ChallengeInstance{
				Status:         challengeInstance.Status,
				InstanceConfig: challengeInstance.InstanceConfig,
			}).
			Error
		if err != nil {
			return nil, errcode.ErrAgentUpdateState.Wrap(err)
		}
	}

	ret := &AgentUpdateState_Output{}
	return ret, nil
}
