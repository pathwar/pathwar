package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestService_OrganizationCreate(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	_, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	tests := []struct {
		name        string
		input       *OrganizationCreate_Input
		expectedErr error
	}{
		{"empty", &OrganizationCreate_Input{}, errcode.ErrMissingInput},
		{"nil", nil, errcode.ErrMissingInput},
		{"reserved name", &OrganizationCreate_Input{Name: "pathwar"}, errcode.ErrReservedName},
		{"my_organization", &OrganizationCreate_Input{Name: "my_organization"}, nil},
		{"my_organization", &OrganizationCreate_Input{Name: "my_organization"}, errcode.ErrCheckOrganizationUniqueName},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := svc.OrganizationCreate(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
		})
	}
}
