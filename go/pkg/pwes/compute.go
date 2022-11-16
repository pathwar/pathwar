package pwes

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type challengeValidation struct {
	validations []*pwdb.Activity
	score       int64
}

//TODO: Don't forget to test the function

func Compute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{FilteringPreset: "validations"})
	if err != nil {
		return err
	}
	activities := res.GetActivities()

	challenges := make(map[int64]*challengeValidation)
	for _, activity := range activities {
		if _, ok := challenges[activity.ChallengeID]; ok {
			challenges[activity.ChallengeID].validations = append(challenges[activity.ChallengeID].validations, activity)
		} else {
			challenges[activity.ChallengeID] = &challengeValidation{[]*pwdb.Activity{activity}, 0}
		}
	}

	// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
	for _, challenge := range challenges {
		nbValidations := len(challenge.validations)
		challenge.score = int64(1/(nbValidations/10+1)*95 + 5)
	}

	teams := make(map[int64]*pwdb.Team)
	for _, activity := range activities {
		if _, ok := teams[activity.TeamID]; !ok {
			teams[activity.TeamID] = activity.Team
			teams[activity.TeamID].Score = challenges[activity.ChallengeID].score
		} else {
			teams[activity.TeamID].Score += challenges[activity.ChallengeID].score
		}
	}

	return nil
}
