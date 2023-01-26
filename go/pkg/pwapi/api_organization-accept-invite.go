package pwapi

import (
	"context"

	"github.com/jinzhu/gorm"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) OrganizationAcceptInvite(ctx context.Context, in *OrganizationAcceptInvite_Input) (*OrganizationAcceptInvite_Output, error) {
	if in == nil || in.OrganizationInviteID == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	OrganizationInviteID, err := pwdb.GetIDBySlugAndKind(svc.db, in.OrganizationInviteID, "organization-invite")
	if err != nil {
		return nil, err
	}

	var organizationInvite pwdb.OrganizationInvite
	err = svc.db.
		Where(pwdb.OrganizationInvite{
			ID:     OrganizationInviteID,
			UserID: userID,
		}).Preload("Organization").
		First(&organizationInvite).
		Error
	if err != nil {
		return nil, pwdb.GormToErrcode(err)
	}

	// check organization status
	var organization pwdb.Organization
	err = svc.db.
		Where(pwdb.Organization{
			ID:             organizationInvite.Organization.ID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		First(&organization).
		Error
	if err != nil {
		return nil, errcode.ErrOrganizationDoesNotExist.Wrap(err)
	}

	// check if the user is already a member of the organization
	var organizationMembership int
	err = svc.db.
		Model(&pwdb.OrganizationMember{}).
		Where(pwdb.OrganizationMember{
			UserID:         userID,
			OrganizationID: organization.ID,
		}).
		Count(&organizationMembership).
		Error
	if err != nil || organizationMembership != 0 {
		return nil, errcode.ErrOrganizationUserAlreadyMember.Wrap(err)
	}

	orgaMember := &pwdb.OrganizationMember{
		UserID:         userID,
		OrganizationID: organization.ID,
	}

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(orgaMember).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}
		err = tx.Delete(&organizationInvite).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}

		activity := pwdb.Activity{
			Kind:                 pwdb.Activity_OrganizationInviteAccept,
			AuthorID:             userID,
			UserID:               organizationInvite.UserID,
			OrganizationID:       organization.ID,
			OrganizationMemberID: orgaMember.ID,
		}

		err = tx.Create(&activity).Error
		if err != nil {
			return pwdb.GormToErrcode(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	ret := OrganizationAcceptInvite_Output{
		OrganizationMember: orgaMember,
	}
	return &ret, nil
}
