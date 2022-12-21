package pwes

import (
	"context"
	"fmt"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func (e EventChallengeSubscriptionValidate) execute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	if apiClient == nil {
		return errcode.ErrMissingInput
	}
	fmt.Println("Input : ", e.SeasonChallenge.ID)
	challenge, err := apiClient.SeasonChallengeGet(ctx, &pwapi.SeasonChallengeGet_Input{SeasonChallengeID: e.SeasonChallenge.ID})
	if err != nil || challenge.Item == nil {
		return errcode.TODO.Wrap(err)
	}
	fmt.Println("Item Slug : ", challenge.Item.Slug)
	return nil
}
