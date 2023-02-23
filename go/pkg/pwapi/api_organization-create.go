package pwapi

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/jinzhu/gorm"

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

	// check for gravatar image
	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%x", md5.Sum([]byte(in.GravatarMail)))

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
		GravatarURL:    gravatarURL,
		// Locale
	}

	// save new organization object in database
	err = svc.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&organization).Error; err != nil {
			return errcode.ErrCreateOrganization.Wrap(err)
		}
		activity := pwdb.Activity{
			Kind:           pwdb.Activity_OrganizationCreation,
			AuthorID:       userID,
			OrganizationID: organization.ID,
		}
		return tx.Create(&activity).Error
	})
	if err != nil {
		return nil, errcode.ErrCreateOrganization.Wrap(err)
	}

	ret := &OrganizationCreate_Output{Organization: &organization}
	return ret, nil
}
