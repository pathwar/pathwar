package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AgentList(ctx context.Context, in *AgentList_Input) (*AgentList_Output, error) {
	if !isAgentContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	return nil, errcode.ErrNotImplemented
}
