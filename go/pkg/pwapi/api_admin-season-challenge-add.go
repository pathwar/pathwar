package pwapi

import (
	"context"

	"gorm.io/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminSeasonChallengeAdd(ctx context.Context, in *AdminSeasonChallengeAdd_Input) (*AdminSeasonChallengeAdd_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	in.ApplyDefaults()
	if in == nil ||
		(in.FlavorID == "" && in.SeasonChallenge.FlavorID == 0) ||
		(in.SeasonID == "" && in.SeasonChallenge.SeasonID == 0) {
		return nil, errcode.ErrMissingInput
	}

	if in.FlavorID != "" && in.SeasonChallenge.FlavorID == 0 {
		var err error
		in.SeasonChallenge.FlavorID, err = pwdb.GetIDBySlugAndKind(svc.db, in.FlavorID, "challenge-flavor")
		if err != nil {
			return nil, err
		}
	}

	if in.SeasonID != "" && in.SeasonChallenge.SeasonID == 0 {
		var err error
		in.SeasonChallenge.SeasonID, err = pwdb.GetIDBySlugAndKind(svc.db, in.SeasonID, "season")
		if err != nil {
			return nil, err
		}
	}

	var seasonChallenge pwdb.SeasonChallenge
	err := svc.db.
		Where(&pwdb.SeasonChallenge{
			SeasonID: in.SeasonChallenge.SeasonID,
			FlavorID: in.SeasonChallenge.FlavorID,
		}).
		First(&seasonChallenge).
		Error
	if err == nil {
		in.SeasonChallenge.ID = seasonChallenge.ID
	} else if err != gorm.ErrRecordNotFound {
		return nil, pwdb.GormToErrcode(err)
	}

	err = svc.db.Save(in.SeasonChallenge).Error
	if err != nil {
		return nil, errcode.ErrSeasonChallengeAdd.Wrap(err)
	}

	out := AdminSeasonChallengeAdd_Output{
		SeasonChallenge: in.SeasonChallenge,
	}
	return &out, nil
}

func (in *AdminSeasonChallengeAdd_Input) ApplyDefaults() {
	if in == nil {
		return
	}
	if in.SeasonChallenge == nil {
		in.SeasonChallenge = &pwdb.SeasonChallenge{}
	}
	if in.SeasonID == "" {
		in.SeasonID = "global"
	}
}
