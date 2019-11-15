package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeGet(ctx context.Context, in *SeasonChallengeGetInput) (*SeasonChallengeGetOutput, error) {
	{ // validation
		if in.SeasonChallengeID == 0 {
			return nil, ErrMissingArgument
		}
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	season, err := seasonFromSeasonChallengeID(e.db, in.SeasonChallengeID)
	if err != nil {
		return nil, ErrInvalidArgument // season challenge is malformed
	}

	team, err := userTeamForSeason(e.db, userID, season.ID)
	if err != nil {
		return nil, ErrInvalidArgument // user does not have team for this season
	}

	var item pwdb.SeasonChallenge
	err = e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.SeasonChallenge{ID: in.SeasonChallengeID}).
		Preload("Season").
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Subscriptions", "team_id = ?", team.ID).
		Preload("Subscriptions.Validations").
		First(&item).
		Error

	switch {
	case err != nil && pwdb.IsRecordNotFoundError(err):
		return nil, ErrInvalidArgument // FIXME: wrap original error
	case err != nil:
		return nil, fmt.Errorf("query season challenge: %w", err)
	}

	ret := SeasonChallengeGetOutput{
		Item: &item,
	}

	return &ret, nil
}
