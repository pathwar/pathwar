package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) TeamList(ctx context.Context, in *TeamList_Input) (*TeamList_Output, error) {
	if in == nil || in.SeasonID == 0 {
		return nil, errcode.ErrMissingInput
	}

	exists, err := seasonIDExists(svc.db, in.SeasonID)
	if err != nil || !exists {
		return nil, errcode.ErrInvalidSeasonID.Wrap(err)
	}

	// query
	var ret TeamList_Output
	err = svc.db.
		Preload("Organization").
		// Preload("Season").
		// Preload("Members").
		// Preload("ChallengeSubscription").
		// Preload("Achievements").
		Where(pwdb.Team{
			SeasonID:       in.SeasonID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, errcode.ErrGetTeams.Wrap(err)
	}

	return &ret, nil
}
