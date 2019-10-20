package pwengine

import (
	"context"
)

func (e *engine) ListChallenges(context.Context, *Void) (*ListChallengesOutput, error) {
	var challenges ListChallengesOutput
	if err := e.db.Set("gorm:auto_preload", true).Find(&challenges.Items).Error; err != nil {
		return nil, err
	}

	return &challenges, nil
}
