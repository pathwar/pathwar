package pwes

import (
	"context"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type challengeValidation struct {
	seasonChallenge *pwdb.SeasonChallenge
	score           int64
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

	*timestamp = *activities[len(activities)-1].CreatedAt

	challengesMap := make(map[int64]*challengeValidation)
	for _, activity := range activities {
		if _, ok := challengesMap[activity.SeasonChallengeID]; ok {
			challengesMap[activity.SeasonChallengeID].seasonChallenge.NbValidations += 1
		} else {
			challengesMap[activity.SeasonChallengeID] = &challengeValidation{&pwdb.SeasonChallenge{ID: activity.SeasonChallengeID, NbValidations: 1}, 0}
		}
	}

	// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
	for _, challenge := range challengesMap {
		challenge.score = 1/(challenge.seasonChallenge.NbValidations/10+1)*95 + 5
	}

	teamsMap := make(map[int64]*pwdb.Team)
	for _, activity := range activities {
		if _, ok := teamsMap[activity.TeamID]; !ok {
			teamsMap[activity.TeamID] = activity.Team
			teamsMap[activity.TeamID].Score = challengesMap[activity.SeasonChallengeID].score
		} else {
			teamsMap[activity.TeamID].Score += challengesMap[activity.SeasonChallengeID].score
		}
	}

	seasonChallenges := []*pwdb.SeasonChallenge{}
	for _, seasonChallenge := range challengesMap {
		seasonChallenges = append(seasonChallenges, seasonChallenge.seasonChallenge)
	}
	_, err = apiClient.AdminUpdateValidations(ctx, &pwapi.AdminUpdateValidations_Input{SeasonChallenge: seasonChallenges})

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
