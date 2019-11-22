package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeList(ctx context.Context, in *SeasonChallengeList_Input) (*SeasonChallengeList_Output, error) {
	if in == nil || in.SeasonID == 0 {
		return nil, ErrMissingArgument
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	exists, err := seasonIDExists(e.db, in.SeasonID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if !exists {
		return nil, ErrInvalidArgument
	}

	team, err := userTeamForSeason(e.db, userID, in.SeasonID)
	if err != nil {
		return nil, ErrInvalidArgument // user does not have team for this season
	}

	var ret SeasonChallengeList_Output
	err = e.db.
		Preload("Season").
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Subscriptions", "team_id = ?", team.ID).
		Preload("Subscriptions.Validations").
		Where(pwdb.SeasonChallenge{SeasonID: in.SeasonID}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch season challenges: %w", err)
	}

	return &ret, nil
}
