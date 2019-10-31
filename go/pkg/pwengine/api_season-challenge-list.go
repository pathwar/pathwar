package pwengine

import (
	"context"
	"fmt"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) SeasonChallengeList(ctx context.Context, in *SeasonChallengeListInput) (*SeasonChallengeListOutput, error) {
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

	var ret SeasonChallengeListOutput
	err := e.db.
		Set("gorm:auto_preload", true).
		Where(pwdb.SeasonChallenge{SeasonID: in.SeasonID}).
		Find(&ret.Items).
		Error
	if err != nil {
		return nil, fmt.Errorf("fetch season challenges: %w", err)
	}

	return &ret, nil
}
