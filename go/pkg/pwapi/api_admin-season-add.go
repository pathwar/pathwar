package pwapi

import (
	"context"

	"gorm.io/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminSeasonAdd(ctx context.Context, in *AdminSeasonAdd_Input) (*AdminSeasonAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.Season == nil {
		return nil, errcode.ErrMissingInput
	}

	// check that the name is not taken
	var seasonCheck pwdb.Season
	err := svc.db.Where(&pwdb.Season{Name: in.Season.Name}).First(&seasonCheck).Error
	if err == nil {
		return nil, errcode.ErrSeasonNameAlreadyExist
	}
	if err != gorm.ErrRecordNotFound {
		return nil, pwdb.GormToErrcode(err)
	}

	err = svc.db.Create(&in.Season).Error
	if err != nil {
		return nil, errcode.ErrSeasonChallengeAdd.Wrap(err)
	}

	out := AdminSeasonAdd_Output{
		Season: in.Season,
	}
	return &out, nil
}
