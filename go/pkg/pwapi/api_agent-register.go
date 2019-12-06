package pwapi

import (
	"context"

	"pathwar.land/go/pkg/errcode"
)

func (svc *service) AgentRegister(context.Context, *AgentRegister_Input) (*AgentRegister_Output, error) {
	// FIXME: check if client is agent
	return nil, errcode.ErrNotImplemented
}
