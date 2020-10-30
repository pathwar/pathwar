package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListOrganizations(ctx context.Context, in *AdminListOrganizations_Input) (*AdminListOrganizations_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var organizations []*pwdb.Organization
	err := svc.db.
		Preload("Teams").
		Preload("Teams.Season").
		Preload("Members").
		Preload("Members.User").
		// Preload("ReceivedWhoswhoAttempts").
		Find(&organizations).Error
	if err != nil {
		return nil, errcode.ErrListOrganizations.Wrap(err)
	}

	out := AdminListOrganizations_Output{Organizations: organizations}
	return &out, nil
}
