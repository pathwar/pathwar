package pwengine

import (
	"context"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeGet(ctx context.Context, in *SeasonChallengeGet_Input) (*SeasonChallengeGet_Output, error) {
	if in == nil || in.SeasonChallengeID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	season, err := seasonFromSeasonChallengeID(e.db, in.SeasonChallengeID)
	if err != nil {
		return nil, errcode.ErrGetSeasonFromSeasonChallenge.Wrap(err)
	}

	team, err := userTeamForSeason(e.db, userID, season.ID)
	if err != nil {
		return nil, errcode.ErrGetUserTeamFromSeason.Wrap(err)
	}

	var item pwdb.SeasonChallenge
	err = e.db.
		Where(pwdb.SeasonChallenge{ID: in.SeasonChallengeID}).
		Preload("Season").
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Subscriptions", "team_id = ?", team.ID).
		Preload("Subscriptions.Validations").
		First(&item).
		Error
	if err != nil {
		return nil, errcode.ErrGetSeasonChallenge.Wrap(err)
	}

	ret := SeasonChallengeGet_Output{Item: &item}
	return &ret, nil
}
