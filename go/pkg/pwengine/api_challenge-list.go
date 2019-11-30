package pwengine

import (
	"context"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ChallengeList(context.Context, *ChallengeList_Input) (*ChallengeList_Output, error) {
	return nil, errcode.ErrDeprecated

	var challenges ChallengeList_Output
	err := e.db.
		Preload("Flavors").
		Find(&challenges.Items).Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	return &challenges, nil
}
