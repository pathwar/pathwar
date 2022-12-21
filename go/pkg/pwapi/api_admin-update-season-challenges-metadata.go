package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminUpdateSeasonChallengesMetadata(ctx context.Context, in *AdminUpdateSeasonChallengesMetadata_Input) (*AdminUpdateSeasonChallengesMetadata_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.SeasonChallenges == nil {
		return nil, errcode.ErrMissingInput
	}

	for _, challenge := range in.SeasonChallenges {
		err := svc.db.Model(challenge).Update("NbValidations", challenge.NbValidations).Error
		if err != nil {
			return nil, err
		}
	}

	return &AdminUpdateSeasonChallengesMetadata_Output{}, nil
}
