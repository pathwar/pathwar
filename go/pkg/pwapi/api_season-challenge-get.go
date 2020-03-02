package pwapi

import (
	"context"

	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func (svc *service) SeasonChallengeGet(ctx context.Context, in *SeasonChallengeGet_Input) (*SeasonChallengeGet_Output, error) {
	if in == nil || in.SeasonChallengeID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	season, err := seasonFromSeasonChallengeID(svc.db, in.SeasonChallengeID)
	if err != nil {
		return nil, errcode.ErrGetSeasonFromSeasonChallenge.Wrap(err)
	}

	team, err := userTeamForSeason(svc.db, userID, season.ID)
	if err != nil {
		return nil, errcode.ErrGetUserTeamFromSeason.Wrap(err)
	}

	var item pwdb.SeasonChallenge
	err = svc.db.
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
