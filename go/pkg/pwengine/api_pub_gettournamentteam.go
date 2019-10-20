package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) GetTournamentTeam(ctx context.Context, in *GetTournamentTeamInput) (*GetTournamentTeamOutput, error) {
	{ // validation
		if in.TournamentTeamID == 0 {
			return nil, ErrMissingArgument
		}
	}

	var item pwdb.TournamentTeam
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.TournamentTeam{ID: in.TournamentTeamID}).
		First(&item).
		Error

	switch {
	case err != nil && pwdb.IsRecordNotFoundError(err):
		return nil, ErrInvalidArgument // FIXME: wrap original error
	case err != nil:
		return nil, fmt.Errorf("fetch tournament team from db: %w", err)
	}

	ret := GetTournamentTeamOutput{
		Item: &item,
	}

	return &ret, nil
}
