package pwengine

import (
	"context"

	"pathwar.land/go/pkg/errcode"
)

func (e *engine) TeamSendInvite(ctx context.Context, in *TeamSendInvite_Input) (*TeamSendInvite_Output, error) {
	if in == nil || in.TeamID == 0 || in.UserID == 0 {
		return nil, errcode.ErrMissingInput
	}

	ret := TeamSendInvite_Output{}
	return &ret, errcode.ErrNotImplemented
}
