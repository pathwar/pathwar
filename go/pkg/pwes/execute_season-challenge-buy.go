package pwes

import (
	"context"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

// TODO: Verify that the user have enough money to buy the challenge
func (e EventSeasonChallengeBuy) execute(ctx context.Context, apiClient *pwapi.HTTPClient, logger *zap.Logger) error {
	if apiClient == nil {
		logger.Debug("missing apiClient in execute event method")
		return errcode.ErrMissingInput
	}
	res, err := apiClient.SeasonChallengeGet(ctx, &pwapi.SeasonChallengeGet_Input{SeasonChallengeID: e.SeasonChallenge.ID})
	challenge := res.GetItem()
	if err != nil || challenge == nil {
		return errcode.TODO.Wrap(err)
	}

	if e.Team.Cash-challenge.Flavor.PurchasePrice < 0 {
		return errcode.ErrNotEnoughCash
	}

	e.Team.Cash -= challenge.Flavor.PurchasePrice
	_, err = apiClient.AdminUpdateTeamsMetadata(ctx, &pwapi.AdminUpdateTeamsMetadata_Input{Teams: []*pwdb.Team{e.Team}})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	return nil
}
