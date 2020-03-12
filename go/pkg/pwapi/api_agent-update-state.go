package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
)

func (svc *service) AgentUpdateState(ctx context.Context, in *AgentUpdateState_Input) (*AgentUpdateState_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	for _, challengeInstance := range in.Instances {
		err := svc.db.Model(&challengeInstance).Update("Status", challengeInstance.Status).Error
		if err != nil {
			return nil, errcode.ErrAgentUpdateState.Wrap(err)
		}
	}

	ret := &AgentUpdateState_Output{}
	return ret, nil
}
