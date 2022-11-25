package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminUpdateValidations(ctx context.Context, in *AdminUpdateValidations_Input) (*AdminUpdateValidations_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.SeasonChallenge == nil {
		return nil, errcode.ErrMissingInput
	}

	// I need Where 1=1 bc of that : https://gorm.io/docs/update.html#block_global_updates
	for _, challenge := range in.SeasonChallenge {
		err := svc.db.Model(&pwdb.SeasonChallenge{}).Where("1 = 1").Update("nb_validations", challenge.NbValidations).Error
		if err != nil {
			return nil, err
		}
	}

	return &AdminUpdateValidations_Output{}, nil
}
