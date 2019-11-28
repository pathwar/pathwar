package pwengine

import (
	"context"
	"math/rand"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) TeamList(ctx context.Context, in *TeamList_Input) (*TeamList_Output, error) {
	if in == nil || in.SeasonID == 0 {
		return nil, errcode.ErrMissingInput
	}

	exists, err := seasonIDExists(e.db, in.SeasonID)
	if err != nil || !exists {
		return nil, errcode.ErrInvalidSeasonID.Wrap(err)
	}

	// query
	var ret TeamList_Output
	err = e.db.
		Preload("Organization").
		//Preload("Season").
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
		return nil, errcode.ErrGetTeams.Wrap(err)
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
