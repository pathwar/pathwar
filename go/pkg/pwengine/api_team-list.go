package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) TeamList(ctx context.Context, in *TeamListInput) (*TeamListOutput, error) {
	{ // validation
		if in.SeasonID == 0 {
			return nil, ErrMissingArgument
		}

		var c int
		err := e.db.
			Table("season").
			Select("id").
			Where(&pwdb.Season{ID: in.SeasonID}).
			Count(&c).
			Error
		if err != nil {
			return nil, fmt.Errorf("fetch season: %w", err)
		}
		if c == 0 {
			return nil, ErrInvalidArgument // invalid in.SeasonID
		}
	}

	var ret TeamListOutput
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.Team{SeasonID: in.SeasonID}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch season organizations from db: %w", err)
	}

	return &ret, nil
}
