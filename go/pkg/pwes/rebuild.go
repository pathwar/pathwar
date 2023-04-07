package pwes

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/errcode"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type ChallengeValidation struct {
	seasonChallenge *pwdb.SeasonChallenge
	score           int64
}

// Rebuild TODO: Rebuild current state from all past events
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

	challengesMap, seasonChallengesID := RebuildNbValidations(activities)

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

	teamsMap, _ := RebuildScore(activities, challengesMap)

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

func RebuildStats(ctx context.Context, apiClient *pwapi.HTTPClient, to *time.Time, seasonID string) (pwapi.AdminSeasonStats_Output, error) {
	if apiClient == nil {
		return pwapi.AdminSeasonStats_Output{}, errcode.ErrMissingInput
	}

	res, err := apiClient.AdminListActivities(ctx, &pwapi.AdminListActivities_Input{FilteringPreset: "validations", To: to})
	if err != nil {
		return pwapi.AdminSeasonStats_Output{}, errcode.TODO.Wrap(err)
	}

	activities := res.GetActivities()
	if len(activities) == 0 {
		return pwapi.AdminSeasonStats_Output{}, errcode.ErrNothingToRebuild
	}

	challengesMap, _ := RebuildNbValidations(activities)
	for _, challenge := range challengesMap {
		challenge.score = computeScore(challenge.seasonChallenge.NbValidations)
	}

	teamsMap, challengesSolvedPerTeam := RebuildScore(activities, challengesMap)
	teams := []*pwdb.Team{}
	for _, team := range teamsMap {
		teams = append(teams, team)
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Score > teams[j].Score
	})

	out := pwapi.AdminSeasonStats_Output{}
	for rank, team := range teams {
		if strconv.FormatInt(team.SeasonID, 10) != seasonID && team.Slug != seasonID {
			continue
		}
		teamPreload, err := apiClient.TeamGet(ctx, &pwapi.TeamGet_Input{TeamID: team.ID})
		if err != nil {
			continue
		}
		for _, member := range teamPreload.GetItem().GetMembers() {
			stat := pwapi.AdminSeasonStats_Output_Stat{
				Rank:             strconv.FormatInt(int64(rank+1), 10),
				Mail:             member.User.Email,
				Name:             member.User.Slug,
				TeamName:         team.Slug[:strings.LastIndex(team.Slug, "@")],
				Score:            strconv.FormatInt(team.Score, 10),
				ChallengesSolved: strconv.FormatInt(challengesSolvedPerTeam[team.ID], 10),
			}
			out.Stats = append(out.Stats, &stat)
		}
	}
	return out, nil
}

func RebuildNbValidations(activities []*pwdb.Activity) (map[int64]*ChallengeValidation, []int64) {
	challengesMap := make(map[int64]*ChallengeValidation)
	var seasonChallengesID []int64
	for _, activity := range activities {
		if _, ok := challengesMap[activity.SeasonChallengeID]; ok {
			challengesMap[activity.SeasonChallengeID].seasonChallenge.NbValidations++
		} else {
			challengesMap[activity.SeasonChallengeID] = &ChallengeValidation{&pwdb.SeasonChallenge{ID: activity.SeasonChallengeID, NbValidations: 1}, 0}
			seasonChallengesID = append(seasonChallengesID, activity.SeasonChallengeID)
		}
	}
	return challengesMap, seasonChallengesID
}

func RebuildScore(activities []*pwdb.Activity, challengesMap map[int64]*ChallengeValidation) (map[int64]*pwdb.Team, map[int64]int64) {
	teamsMap := make(map[int64]*pwdb.Team)
	challengesSolvedPerTeam := make(map[int64]int64)
	for _, activity := range activities {
		if _, ok := teamsMap[activity.TeamID]; !ok {
			teamsMap[activity.TeamID] = activity.Team
			teamsMap[activity.TeamID].Score = 0
			challengesSolvedPerTeam[activity.TeamID] = 0
		}
		teamsMap[activity.TeamID].Score += challengesMap[activity.SeasonChallengeID].score
		challengesSolvedPerTeam[activity.TeamID]++
	}
	return teamsMap, challengesSolvedPerTeam
}
