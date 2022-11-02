package pwapi

import (
	"context"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
)

func TestService_SeasonList(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	tests := []struct {
		name              string
		input             *SeasonList_Input
		expectedErr       error
		ExpectedStatus    pwdb.Season_Status
		ExpectedGlobal    bool
		ExpectedIsTesting bool
	}{
		{"all-seasons", &SeasonList_Input{}, nil, pwdb.Season_Started, true, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.SeasonList(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			for _, season := range ret.Items {
				assert.Equal(t, test.ExpectedStatus, season.Status)
				assert.Equal(t, test.ExpectedGlobal, season.IsGlobal)
				assert.Equal(t, test.ExpectedIsTesting, season.IsTesting)
			}
		})
	}
}
