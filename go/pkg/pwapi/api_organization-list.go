package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) OrganizationList(context.Context, *OrganizationList_Input) (*OrganizationList_Output, error) {
	var organizations OrganizationList_Output
	err := svc.db.
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
