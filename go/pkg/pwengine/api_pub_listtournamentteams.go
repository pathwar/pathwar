package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ListTournamentTeams(ctx context.Context, in *ListTournamentTeamsInput) (*pwdb.TournamentTeamList, error) {
	{ // validation
		if in.TournamentID == "" {
			return nil, ErrMissingArgument
		}

		var c int
		err := e.db.
			Table("tournament").
			Select("id").
			Where("id = ?", in.TournamentID).
			Count(&c).
			Error
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tournament: %w", err)
		}
		if c == 0 {
			return nil, ErrInvalidArgument // invalid in.TournamentID
		}
	}

	var ret pwdb.TournamentTeamList
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.TournamentTeam{TournamentID: in.TournamentID}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tournament teams from db: %w", err)
	}

	return &ret, nil
}
