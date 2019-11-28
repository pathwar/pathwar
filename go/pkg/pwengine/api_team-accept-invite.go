package pwengine

import (
	"context"

	"pathwar.land/go/pkg/errcode"
)

func (e *engine) TeamAcceptInvite(ctx context.Context, in *TeamAcceptInvite_Input) (*TeamAcceptInvite_Output, error) {
	if in == nil || in.TeamMemberID == 0 {
		return nil, errcode.ErrMissingInput
	}

	ret := TeamAcceptInvite_Output{}
	return &ret, errcode.ErrNotImplemented
}
