package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) OrganizationSendInvite(ctx context.Context, in *OrganizationSendInvite_Input) (*OrganizationSendInvite_Output, error) {
	if in == nil || in.OrganizationID == "" || in.UserID == "" {
		return nil, errcode.ErrMissingInput
	}

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated.Wrap(err)
	}

	organizationID, err := pwdb.GetIDBySlugAndKind(svc.db, in.OrganizationID, "organization")
	if err != nil {
		return nil, errcode.ErrGetOrganization.Wrap(err)
	}

	inviteUserID, err := pwdb.GetIDBySlugAndKind(svc.db, in.UserID, "user")
	if err != nil {
		return nil, errcode.ErrGetUser.Wrap(err)
	}

	//check organization status
	var organization pwdb.Organization
	err = svc.db.
		Where(pwdb.Organization{
			ID:             organizationID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		First(&organization).
		Error
	if err != nil {
		return nil, errcode.ErrOrganizationDoesNotExist.Wrap(err)
	}

	// check that the user is owner of the organization
	var organizationOwner pwdb.OrganizationMember
	err = svc.db.
		Where(pwdb.OrganizationMember{
			UserID:         userID,
			OrganizationID: organizationID,
			Role:           pwdb.OrganizationMember_Owner,
		}).
		First(&organizationOwner).
		Error
	if err != nil {
		return nil, errcode.ErrNotOrganizationOwner.Wrap(err)
	}

	// check if invited user is not already a member of the organization
	var organizationMember pwdb.OrganizationMember
	err = svc.db.
		Where(pwdb.OrganizationMember{
			UserID:         inviteUserID,
			OrganizationID: organizationID,
		}).
		First(&organizationMember).
		Error
	if err != nil {
		return nil, errcode.ErrOrganizationUserAlreadyMember.Wrap(err)
	}

	return nil, errcode.ErrNotImplemented
}
