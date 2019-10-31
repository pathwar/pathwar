package pwengine

import "context"

func (e *engine) ChallengeSubscriptionValidate(ctx context.Context, in *ChallengeSubscriptionValidateInput) (*ChallengeSubscriptionValidateOutput, error) {
	{ // validation
		if in.ChallengeSubscriptionID == 0 {
			return nil, ErrMissingArgument
		}
	}

	ret := ChallengeSubscriptionValidateOutput{}

	return &ret, nil
}
