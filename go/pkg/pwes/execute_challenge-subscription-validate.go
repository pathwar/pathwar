package pwes

import (
	"context"
	"fmt"
	"strconv"

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
	teams := []*pwdb.Team{}
	if oldScore != newScore {
		res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{SeasonChallengeID: strconv.Itoa(int(e.SeasonChallenge.ID)), FilteringPreset: "validations"})
		validations := res.GetActivities()
		if err != nil {
			return errcode.TODO.Wrap(err)
		}
		for _, validation := range validations {
			teams = append(teams, validation.Team)
		}
	} else {
		teams = append(teams, e.Team)
	}
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
