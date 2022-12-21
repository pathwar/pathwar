package pwes

import (
	"context"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/errcode"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type challengeValidation struct {
	seasonChallenge *pwdb.SeasonChallenge
	score           int64
}

func Rebuild(ctx context.Context, apiClient *pwapi.HTTPClient, opts Opts) error {
	if apiClient == nil {
		return errcode.ErrMissingInput
	}

	if opts.WithoutScore {
		return errcode.ErrNothingToRebuild
	}

	from, _ := time.Parse(TimeLayout, opts.From)
	to, _ := time.Parse(TimeLayout, opts.To)
	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{Since: &from, FilteringPreset: "validations", To: &to})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	activities := res.GetActivities()
	if len(activities) == 0 {
		return errcode.ErrNothingToRebuild
	}

	challengesMap := make(map[int64]*challengeValidation)
	var seasonChallengesID []int64
	for _, activity := range activities {
		if _, ok := challengesMap[activity.SeasonChallengeID]; ok {
			challengesMap[activity.SeasonChallengeID].seasonChallenge.NbValidations++
		} else {
			challengesMap[activity.SeasonChallengeID] = &challengeValidation{&pwdb.SeasonChallenge{ID: activity.SeasonChallengeID, NbValidations: 1}, 0}
			seasonChallengesID = append(seasonChallengesID, activity.SeasonChallengeID)
		}
	}

	listSeasonChallenges, err := apiClient.AdminListSeasonChallenges(ctx, &pwapi.AdminListSeasonChallenges_Input{Id: seasonChallengesID})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	seasonChallenges := listSeasonChallenges.GetSeasonChallenge()
	for _, seasonChallenge := range seasonChallenges {
		seasonChallenge.NbValidations = challengesMap[seasonChallenge.ID].seasonChallenge.NbValidations
	}

	// TODO: Apply a better function: compute score : 1 / (x/10 + 1) * 95 + 5
	for _, challenge := range challengesMap {
		challenge.score = computeScore(challenge.seasonChallenge.NbValidations)
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

	_, err = apiClient.AdminUpdateSeasonChallengesMetadata(ctx, &pwapi.AdminUpdateSeasonChallengesMetadata_Input{SeasonChallenges: seasonChallenges})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	teams := []*pwdb.Team{}
	for _, team := range teamsMap {
		teams = append(teams, team)
	}
	_, err = apiClient.AdminUpdateTeamsMetadata(ctx, &pwapi.AdminUpdateTeamsMetadata_Input{Teams: teams})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}
