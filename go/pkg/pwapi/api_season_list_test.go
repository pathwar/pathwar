package pwapi

import (
	"context"
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

	seasons := testingSeasons(t, svc).Items

	tests := []struct {
		name            string
		input           *SeasonList_Input
		expectedErr     error
		expectedSeasons int
	}{
		{"all-seasons", &SeasonList_Input{}, nil, len(seasons)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.SeasonList(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.expectedSeasons, len(ret.Seasons))
		})
	}
}
