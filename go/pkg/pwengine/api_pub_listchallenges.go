package pwengine

import (
	"context"
	"fmt"
)

func (e *engine) ListChallenges(context.Context, *Void) (*ListChallengesOutput, error) {
	var challenges ListChallengesOutput
	err := e.db.
		Set("gorm:auto_preload", true). // FIXME: explicit preloading
		Find(&challenges.Items).Error
	if err != nil {
		return nil, fmt.Errorf("query challenges: %w", err)
	}

	return &challenges, nil
}
