package pwapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"testing"
)

func TestService_OrganizationList(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	tests := []struct {
		name          string
		input         *OrganizationList_Input
		expectedErr   error
		expectedTeams int
	}{
		{"list", &OrganizationList_Input{}, nil, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.OrganizationList(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, test.expectedTeams, len(ret.Items))
		})
	}
}
