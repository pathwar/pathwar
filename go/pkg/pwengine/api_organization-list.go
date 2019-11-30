package pwengine

import (
	"context"

	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) OrganizationList(context.Context, *OrganizationList_Input) (*OrganizationList_Output, error) {
	var organizations OrganizationList_Output
	err := e.db.
		Preload("Teams").
		// Preload("OrganizationMembers").
		Where(pwdb.Organization{
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		Find(&organizations.Items).
		Error
	if err != nil {
		return nil, errcode.ErrFindOrganizations.Wrap(err)
	}

	return &organizations, nil
}
