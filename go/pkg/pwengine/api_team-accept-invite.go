package pwengine

import (
	"context"
)

func (e *engine) TeamAcceptInvite(ctx context.Context, in *TeamAcceptInvite_Input) (*TeamAcceptInvite_Output, error) {
	// validation
	if in.TeamMemberID == 0 {
		return nil, ErrMissingArgument
	}

	ret := TeamAcceptInvite_Output{}
	return &ret, nil
}
