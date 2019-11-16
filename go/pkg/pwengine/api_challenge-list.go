package pwengine

import (
	"context"
	"fmt"
)

func (e *engine) ChallengeList(context.Context, *ChallengeList_Input) (*ChallengeList_Output, error) {
	return nil, fmt.Errorf("admin call (deprecated)")

	var challenges ChallengeList_Output
	err := e.db.
		Preload("Flavors").
		Find(&challenges.Items).Error
	if err != nil {
		return nil, fmt.Errorf("query challenges: %w", err)
	}

	return &challenges, nil
}
