package pwengine

import (
	"context"
	"fmt"
	"math/rand"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) TeamGet(ctx context.Context, in *TeamGet_Input) (*TeamGet_Output, error) {
	if in == nil || in.TeamID == 0 {
		return nil, ErrMissingArgument
	}

	var item pwdb.Team
	err := e.db.
		Preload("Season").
		Preload("Organization").
		Preload("Members").                // only if member of the team or if admin
		Preload("ChallengeSubscriptions"). // only if member of the team or if admin
		Preload("Achievements").
		Where(pwdb.Team{
			ID:             in.TeamID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		First(&item).
		Error

	switch {
	case err != nil && pwdb.IsRecordNotFoundError(err):
		return nil, ErrInvalidArgument // FIXME: wrap original error
	case err != nil:
		return nil, fmt.Errorf("fetch team from db: %w", err)
	}

	ret := TeamGet_Output{
		Item: &item,
	}

	// tmp: fake data
	// FIXME: use real data
	item.GoldMedals = int64(rand.Intn(3))
	item.SilverMedals = int64(rand.Intn(3))
	item.BronzeMedals = int64(rand.Intn(4))
	item.Score = int64(rand.Intn(100))
	item.Cash = int64(rand.Intn(100))
	item.NbAchievements = int64(rand.Intn(10))

	return &ret, nil
}
