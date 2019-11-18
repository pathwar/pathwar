package pwengine

import (
	"context"
)

func (e *engine) TeamSendInvite(ctx context.Context, in *TeamSendInvite_Input) (*TeamSendInvite_Output, error) {
	// validation
	if in.TeamID == 0 {
		return nil, ErrMissingArgument
	}
	if in.UserID == 0 {
		return nil, ErrMissingArgument
	}

	ret := TeamSendInvite_Output{}
	return &ret, nil
}
