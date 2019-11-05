package pwengine

import (
	"context"
	"fmt"
)

func (e *engine) ChallengeList(context.Context, *ChallengeListInput) (*ChallengeListOutput, error) {
	return nil, fmt.Errorf("admin call (deprecated)")

	var challenges ChallengeListOutput
	err := e.db.
		Preload("Flavors").
		Find(&challenges.Items).Error
	if err != nil {
		return nil, fmt.Errorf("query challenges: %w", err)
	}

	return &challenges, nil
}
