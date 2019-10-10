package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ListChallenges(context.Context, *Void) (*pwdb.ChallengeList, error) {
	var challenges pwdb.ChallengeList
	if err := e.db.Set("gorm:auto_preload", true).Find(&challenges.Items).Error; err != nil {
		return nil, err
	}

	return &challenges, nil
}
