package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) TeamAcceptInvite(ctx context.Context, in *TeamAcceptInvite_Input) (*TeamAcceptInvite_Output, error) {
	if in == nil || in.TeamMemberID == 0 {
		return nil, errcode.ErrMissingInput
	}

	ret := TeamAcceptInvite_Output{}
	return &ret, errcode.ErrNotImplemented
}
