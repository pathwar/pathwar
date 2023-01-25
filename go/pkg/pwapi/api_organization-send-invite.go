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

	return nil, errcode.ErrNotImplemented
}
