package pwes

import (
	"context"
	"fmt"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type challengeValidation struct {
	validations []*pwdb.Activity
}

func Compute(ctx context.Context, apiClient *pwapi.HTTPClient) error {
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{FilteringPreset: "validations"})
	if err != nil {
		return err
	}
	activities := res.GetActivities()

	//TODO: Perhaps a better way to split validations per challenges
	challenges := make(map[int64]challengeValidation)
	for _, activity := range activities {
		if _, ok := challenges[activity.ChallengeID]; ok {
			challenges[activity.ChallengeID] = challengeValidation{append(challenges[activity.ChallengeID].validations, activity)}
		} else {
			challenges[activity.ChallengeID] = challengeValidation{[]*pwdb.Activity{activity}}
		}
	}

	fmt.Println(challenges)
	return nil
}
