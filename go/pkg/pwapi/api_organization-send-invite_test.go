package pwapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func TestService_OrganizationSendInvite(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	db := testingSvcDB(t, svc)
	ctx := testingSetContextToken(context.Background(), t)
	_, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	ctx2 := testingSetContextToken2(context.Background(), t)
	session2, err := svc.UserGetSession(ctx2, nil)
	require.NoError(t, err)

	var organization1 pwdb.Organization
	ret, err := svc.OrganizationCreate(ctx, &OrganizationCreate_Input{
		Name: "Test1Organization",
	})
	require.NoError(t, err)
	err = db.Model(pwdb.Organization{}).Where(&pwdb.Organization{ID: ret.Organization.ID}).Update(&pwdb.Organization{DeletionStatus: pwdb.DeletionStatus_Anonymized}).First(&organization1).Error
	require.NoError(t, err)
	ret, err = svc.OrganizationCreate(ctx, &OrganizationCreate_Input{
		Name: "Test2Organization",
	})
	require.NoError(t, err)
	organization2 := ret.Organization
	ret, err = svc.OrganizationCreate(ctx2, &OrganizationCreate_Input{
		Name: "Test4Organization",
	})
	require.NoError(t, err)
	organization3 := ret.Organization

	tests := []struct {
		name        string
		input       *OrganizationSendInvite_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &OrganizationSendInvite_Input{}, errcode.ErrMissingInput},
		{"invalid-organization-id", &OrganizationSendInvite_Input{OrganizationID: "-1", UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrNoSuchSlug},
		{"only-organization-id", &OrganizationSendInvite_Input{OrganizationID: fmt.Sprint(organization1.ID)}, errcode.ErrMissingInput},
		{"only-user-id", &OrganizationSendInvite_Input{UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrMissingInput},
		{"invalid-user-id", &OrganizationSendInvite_Input{OrganizationID: fmt.Sprint(organization1.ID), UserID: "-1"}, errcode.ErrNoSuchSlug},
		{"deleted-organization", &OrganizationSendInvite_Input{OrganizationID: fmt.Sprint(organization1.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrOrganizationDoesNotExist},
		{"valid", &OrganizationSendInvite_Input{OrganizationID: fmt.Sprint(organization2.ID), UserID: fmt.Sprint(session2.User.ID)}, nil},
		{"already-invited", &OrganizationSendInvite_Input{OrganizationID: fmt.Sprint(organization2.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrAlreadyInvitedInOrganization},
		{"is-owner", &OrganizationSendInvite_Input{OrganizationID: fmt.Sprint(organization3.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrNotOrganizationOwner},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := svc.OrganizationSendInvite(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}
		})
	}
}
