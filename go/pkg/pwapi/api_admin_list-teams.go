package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListTeams(ctx context.Context, in *AdminListTeams_Input) (*AdminListTeams_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var teams []*pwdb.Team
	err := svc.db.
		// Preload("").
		Find(&teams).Error
	if err != nil {
		return nil, errcode.ErrListTeams.Wrap(err)
	}

	out := AdminListTeams_Output{Teams: teams}
	return &out, nil
}
