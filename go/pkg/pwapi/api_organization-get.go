package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) OrganizationGet(ctx context.Context, in *OrganizationGet_Input) (*OrganizationGet_Output, error) {
	if in == nil || in.OrganizationID == 0 {
		return nil, errcode.ErrMissingInput
	}

	var item pwdb.Organization
	err := svc.db.
		Preload("Members").
		Where(pwdb.Organization{
			ID:             in.OrganizationID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		First(&item).
		Error
	if err != nil {
		return nil, errcode.ErrGetOrganization.Wrap(err)
	}

	ret := OrganizationGet_Output{Item: &item}

	return &ret, nil
}
