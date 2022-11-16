package pwapi

import (
	"context"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminSetTeams(ctx context.Context, in *AdminSetTeams_Input) (*AdminSetTeams_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil || in.Teams == nil {
		return nil, errcode.ErrMissingInput
	}

	for _, team := range in.Teams {
		err := svc.db.Model(&pwdb.Team{}).Update(team).Error
		if err != nil {
			return nil, err
		}
	}

	return &AdminSetTeams_Output{}, nil
}
