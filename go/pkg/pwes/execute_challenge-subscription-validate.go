package pwes

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

// TODO: Add a better way to get challenges validations, maybe a challenge could be closed without be validated
func (e *EventChallengeSubscriptionValidate) execute(ctx context.Context, apiClient *pwapi.HTTPClient, logger *zap.Logger) error {
	if apiClient == nil {
		logger.Debug("missing apiClient in execute event method")
		return errcode.ErrMissingInput
	}

	oldScore := computeScore(e.SeasonChallenge.NbValidations)
	e.SeasonChallenge.NbValidations++
	newScore := computeScore(e.SeasonChallenge.NbValidations)

	teams := []*pwdb.Team{}
	if oldScore != newScore {
		diffScore := oldScore - newScore
		res, err := apiClient.AdminListChallengeSubscriptions(ctx, &pwapi.AdminListChallengeSubscriptions_Input{SeasonChallengeID: strconv.Itoa(int(e.SeasonChallenge.ID)), FilteringPreset: "closed"})
		validations := res.GetSubscriptions()
		if err != nil {
			return errcode.TODO.Wrap(err)
		}
		for _, validation := range validations {
			if validation.TeamID == e.Team.ID {
				continue
			}
			validation.Team.Score -= diffScore
			teams = append(teams, validation.Team)
		}
	}
	e.Team.Score += newScore
	teams = append(teams, e.Team)

	_, err := apiClient.AdminUpdateSeasonChallengesMetadata(ctx, &pwapi.AdminUpdateSeasonChallengesMetadata_Input{SeasonChallenges: []*pwdb.SeasonChallenge{e.SeasonChallenge}})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	_, err = apiClient.AdminUpdateTeamsMetadata(ctx, &pwapi.AdminUpdateTeamsMetadata_Input{Teams: teams})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	return nil
}
