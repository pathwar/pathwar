package pwengine

import (
	"context"
	"fmt"

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
		Where(pwdb.Team{
			SeasonID:       in.SeasonID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch teams: %w", err)
	}

	return &ret, nil
}
