package pwengine

import (
	"context"
	"fmt"
	"math/rand"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) TeamList(ctx context.Context, in *TeamList_Input) (*TeamList_Output, error) {
	// validation
	if in.SeasonID == 0 {
		return nil, ErrMissingArgument
	}
	exists, err := seasonIDExists(e.db, in.SeasonID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if !exists {
		return nil, ErrInvalidArgument
	}

	// query
	var ret TeamList_Output
	err = e.db.
		//Preload("Season").
		Preload("Organization").
		//Preload("Members").
		//Preload("ChallengeSubscription").
		//Preload("Achievements").
		Where(pwdb.Team{
			SeasonID:       in.SeasonID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch teams: %w", err)
	}

	// add fake data
	// FIXME: use real data instead
	for _, team := range ret.Items {
		team.GoldMedals = int64(rand.Intn(3))
		team.SilverMedals = int64(rand.Intn(3))
		team.BronzeMedals = int64(rand.Intn(4))
		team.Score = int64(rand.Intn(100))
		team.Cash = int64(rand.Intn(100))
		team.NbAchievements = int64(rand.Intn(10))
	}

	return &ret, nil
}
