package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) ChallengeList(context.Context, *ChallengeList_Input) (*ChallengeList_Output, error) {
	return nil, errcode.ErrDeprecated

	/*
		var challenges ChallengeList_Output
		err := svc.db.
			Preload("Flavors").
			Find(&challenges.Items).Error
		if err != nil {
			return nil, pwdb.GormToErrcode(err)
		}

		return &challenges, nil
	*/
}
