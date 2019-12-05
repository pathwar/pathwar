package pwapi

import (
	"context"

	"pathwar.land/go/pkg/errcode"
)

func (svc *service) AgentListInstances(context.Context, *AgentListInstances_Input) (*AgentListInstances_Output, error) {
	// FIXME: check if client is agent
	return nil, errcode.ErrNotImplemented
}
