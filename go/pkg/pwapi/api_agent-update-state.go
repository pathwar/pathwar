package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
)

func (svc *service) AgentUpdateState(context.Context, *AgentUpdateState_Input) (*AgentUpdateState_Output, error) {
	// FIXME: check if client is agent
	return nil, errcode.ErrNotImplemented
}
