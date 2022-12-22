package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminUpdateTeamsMetadata(ctx context.Context, in *AdminUpdateTeamsMetadata_Input) (*AdminUpdateTeamsMetadata_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.Teams == nil {
		return nil, errcode.ErrMissingInput
	}

	for _, team := range in.Teams {
		updates := pwdb.Team{Score: team.Score}
		err := svc.db.Model(team).Update(&updates).Error
		if err != nil {
			return nil, err
		}
	}

	return &AdminUpdateTeamsMetadata_Output{}, nil
}
