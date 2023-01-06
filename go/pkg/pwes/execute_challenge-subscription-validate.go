package pwes

import (
	"context"
	"go.uber.org/zap"
	"strconv"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

// TODO: Add a better way to get challenges validations, maybe a challenge could be closed without be validated
func (e EventChallengeSubscriptionValidate) execute(ctx context.Context, apiClient *pwapi.HTTPClient, logger *zap.Logger) error {
	if apiClient == nil {
		return errcode.ErrMissingInput
	}
	res, err := apiClient.SeasonChallengeGet(ctx, &pwapi.SeasonChallengeGet_Input{SeasonChallengeID: e.SeasonChallenge.ID})
	challenge := res.GetItem()
	if err != nil || challenge == nil {
		return errcode.TODO.Wrap(err)
	}

	oldScore := computeScore(challenge.NbValidations)
	challenge.NbValidations++
	newScore := computeScore(challenge.NbValidations)

	teams := []*pwdb.Team{}
	if oldScore != newScore {
		diffScore := oldScore - newScore
		res, err := apiClient.AdminListChallengeSubscriptions(ctx, &pwapi.AdminListChallengeSubscriptions_Input{SeasonChallengeID: strconv.Itoa(int(e.SeasonChallenge.ID)), FilteringPreset: "closed"})
		validations := res.GetSubscriptions()
		if err != nil {
			return errcode.TODO.Wrap(err)
		}
		for _, validation := range validations {
			validation.Team.Score -= diffScore
			teams = append(teams, validation.Team)
		}
	}
	e.Team.Score += newScore
	e.Team.Cash += challenge.Flavor.ValidationReward
	teams = append(teams, e.Team)

	_, err = apiClient.AdminUpdateSeasonChallengesMetadata(ctx, &pwapi.AdminUpdateSeasonChallengesMetadata_Input{SeasonChallenges: []*pwdb.SeasonChallenge{challenge}})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	_, err = apiClient.AdminUpdateTeamsMetadata(ctx, &pwapi.AdminUpdateTeamsMetadata_Input{Teams: teams})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	return nil
}
