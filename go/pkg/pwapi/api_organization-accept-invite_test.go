package pwapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestService_OrganizationAcceptInvite(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)
	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	ctx2 := testingSetContextToken2(context.Background(), t)
	_, err = svc.UserGetSession(ctx2, nil)
	require.NoError(t, err)

	ret, err := svc.OrganizationCreate(ctx2, &OrganizationCreate_Input{
		Name: "Test1Organization",
	})
	organization1 := ret.Organization
	ret2, err := svc.OrganizationSendInvite(ctx2, &OrganizationSendInvite_Input{
		OrganizationID: fmt.Sprint(organization1.ID),
		UserID:         fmt.Sprint(session.User.ID),
	})
	require.NoError(t, err)
	organizationInvite1 := ret2.OrganizationInvite

	tests := []struct {
		name        string
		input       *OrganizationAcceptInvite_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &OrganizationAcceptInvite_Input{}, errcode.ErrMissingInput},
		{"invalid-invite", &OrganizationAcceptInvite_Input{OrganizationInviteID: "invalid"}, errcode.ErrNoSuchSlug},
		{"valid", &OrganizationAcceptInvite_Input{OrganizationInviteID: fmt.Sprint(organizationInvite1.ID)}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := svc.OrganizationAcceptInvite(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}
		})
	}

}
