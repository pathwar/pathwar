package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
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

	// load seasonRules
	seasonRules := NewSeasonRules()
	err = seasonRules.ParseSeasonRulesString([]byte(in.Season.RulesBundle))
	if err != nil {
		return nil, errcode.ErrParseSeasonRule.Wrap(err)
	}

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&in.Season).Error
		if err != nil {
			return errcode.ErrSeasonChallengeAdd.Wrap(err)
		}

		if !seasonRules.IsStarted() {
			activity := pwdb.Activity{
				Kind:      pwdb.Activity_SeasonOpen,
				CreatedAt: &seasonRules.StartDatetime,
				SeasonID:  in.Season.ID,
			}
			err = tx.Create(&activity).Error
			if err != nil {
				return errcode.ErrSeasonChallengeAdd.Wrap(err)
			}
		}

		if !seasonRules.EndDatetime.IsZero() {
			activity := pwdb.Activity{
				Kind:      pwdb.Activity_SeasonOpen,
				CreatedAt: &seasonRules.EndDatetime,
				SeasonID:  in.Season.ID,
			}
			err = tx.Create(&activity).Error
			if err != nil {
				return errcode.ErrSeasonChallengeAdd.Wrap(err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	out := AdminSeasonAdd_Output{
		Season: in.Season,
	}
	return &out, nil
}
