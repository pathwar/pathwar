package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) OrganizationCreate(ctx context.Context, in *OrganizationCreate_Input) (*OrganizationCreate_Output, error) {
	if in == nil || in.Name == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	in.Name = normalizeName(in.Name)
	if isReservedName(in.Name) {
		return nil, errcode.ErrReservedName
	}

	// check for existing organization with that name
	var count int
	err = svc.db.Model(pwdb.Organization{}).Where(pwdb.Organization{Name: in.Name}).Count(&count).Error
	if err != nil || count != 0 {
		return nil, errcode.ErrCheckOrganizationUniqueName.Wrap(err)
	}

	// create new organization
	organization := pwdb.Organization{
		Name: in.Name,
		Members: []*pwdb.OrganizationMember{
			{
				UserID: userID,
				Role:   pwdb.OrganizationMember_Owner,
			},
		},
		DeletionStatus: pwdb.DeletionStatus_Active,
		// GravatarURL
		// Locale
	}
	err = svc.db.Create(&organization).Error
	if err != nil {
		return nil, errcode.ErrCreateOrganization.Wrap(err)
	}

	ret := &OrganizationCreate_Output{Organization: &organization}
	return ret, nil
}
