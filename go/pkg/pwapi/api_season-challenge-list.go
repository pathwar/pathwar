package pwapi

import (
	"context"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (svc *service) SeasonChallengeList(ctx context.Context, in *SeasonChallengeList_Input) (*SeasonChallengeList_Output, error) {
	if in == nil || in.SeasonID == 0 {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	exists, err := seasonIDExists(svc.db, in.SeasonID)
	if err != nil || !exists {
		return nil, errcode.ErrInvalidSeasonID.Wrap(err)
	}

	team, err := userTeamForSeason(svc.db, userID, in.SeasonID)
	if err != nil {
		return nil, errcode.ErrUserHasNoTeamForSeason.Wrap(err)
	}

	var ret SeasonChallengeList_Output
	err = svc.db.
		Preload("Season").
		Preload("Flavor").
		Preload("Flavor.Challenge").
		Preload("Subscriptions", "team_id = ?", team.ID).
		Preload("Subscriptions.Validations").
		Where(pwdb.SeasonChallenge{SeasonID: in.SeasonID}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, errcode.ErrGetSeasonChallenges.Wrap(err)
	}

	return &ret, nil
}
