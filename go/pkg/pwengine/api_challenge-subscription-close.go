package pwengine

import "context"

func (e *engine) ChallengeSubscriptionClose(ctx context.Context, in *ChallengeSubscriptionCloseInput) (*ChallengeSubscriptionCloseOutput, error) {
	{ // validation
		if in.ChallengeSubscriptionID == 0 {
			return nil, ErrMissingArgument
		}
	}

	ret := ChallengeSubscriptionCloseOutput{}

	return &ret, nil
}
