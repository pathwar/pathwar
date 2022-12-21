package pwes

import (
	"context"
	"fmt"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func (e EventChallengeSubscriptionValidate) execute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	if apiClient == nil {
		return errcode.ErrMissingInput
	}
	fmt.Println("Input : ", e.SeasonChallenge.ID)
	res, err := apiClient.SeasonChallengeGet(ctx, &pwapi.SeasonChallengeGet_Input{SeasonChallengeID: e.SeasonChallenge.ID})
	challenge := res.GetItem()
	if err != nil || challenge == nil {
		return errcode.TODO.Wrap(err)
	}

	oldScore := computeScore(challenge.NbValidations)
	newScore := computeScore(challenge.NbValidations + 1)

	teamsMap := make(map[int64]*pwdb.Team)
	if oldScore != newScore {
		return nil
	}
	var _ = teamsMap
	return nil
}
