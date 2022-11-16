package pwes

import (
	"context"
	"fmt"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type challengeValidation struct {
	validations []*pwdb.Activity
	score       int
}

//TODO: Don't forget to test the function

func Compute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{FilteringPreset: "validations"})
	if err != nil {
		return err
	}
	activities := res.GetActivities()

	//TODO: Perhaps a better way to split validations per challenges
	challenges := make(map[int64]challengeValidation)
	for _, activity := range activities {
		challenges[activity.ChallengeID] = challengeValidation{append(challenges[activity.ChallengeID].validations, activity), 0}
	}

	// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
	for _, challenge := range challenges {
		nbValidations := len(challenge.validations)
		challenge.score = 1/(nbValidations/10+1)*95 + 5
	}

	fmt.Println(challenges)
	return nil
}
