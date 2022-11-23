package pwes

import (
	"context"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type challengeValidation struct {
	validations []*pwdb.Activity
	score       int64
}

func Compute(ctx context.Context, apiClient *pwapi.HTTPClient, timestamp *time.Time) error {
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{Since: timestamp, FilteringPreset: "validations"})
	if err != nil {
		return err
	}

	activities := res.GetActivities()
	if len(activities) == 0 {
		return nil
	}
	
	timestamp = activities[len(activities)-1].CreatedAt

	challengesMap := make(map[int64]*challengeValidation)
	for _, activity := range activities {
		if _, ok := challengesMap[activity.ChallengeID]; ok {
			challengesMap[activity.ChallengeID].validations = append(challengesMap[activity.ChallengeID].validations, activity)
		} else {
			challengesMap[activity.ChallengeID] = &challengeValidation{[]*pwdb.Activity{activity}, 0}
		}
	}

	// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
	for _, challenge := range challengesMap {
		nbValidations := len(challenge.validations)
		challenge.score = int64(1/(nbValidations/10+1)*95 + 5)
	}

	teamsMap := make(map[int64]*pwdb.Team)
	for _, activity := range activities {
		if _, ok := teamsMap[activity.TeamID]; !ok {
			teamsMap[activity.TeamID] = activity.Team
			teamsMap[activity.TeamID].Score = challengesMap[activity.ChallengeID].score
		} else {
			teamsMap[activity.TeamID].Score += challengesMap[activity.ChallengeID].score
		}
	}

	teams := []*pwdb.Team{}
	for _, team := range teamsMap {
		teams = append(teams, team)
	}

	_, err = apiClient.AdminSetTeams(ctx, &pwapi.AdminSetTeams_Input{Teams: teams})
	if err != nil {
		return err
	}

	return nil
}
