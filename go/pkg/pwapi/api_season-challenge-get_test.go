package pwapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/go/v2/internal/testutil"
	"pathwar.land/go/v2/pkg/errcode"
)

func TestSvc_SeasonChallengeGet(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	seasonChallenges := map[string]int64{}
	for _, seasonChallenge := range testingSeasonChallenges(t, svc).Items {
		key := fmt.Sprintf("%s/%s", seasonChallenge.Season.Name, seasonChallenge.Flavor.Challenge.Name)
		seasonChallenges[key] = seasonChallenge.ID
	}

	var tests = []struct {
		name                  string
		input                 *SeasonChallengeGet_Input
		expectedErr           error
		expectedSeasonName    string
		expectedChallengeName string
	}{
		{"empty", &SeasonChallengeGet_Input{}, errcode.ErrMissingInput, "", ""},
		{"unknown-season-id", &SeasonChallengeGet_Input{SeasonChallengeID: -42}, errcode.ErrGetSeasonFromSeasonChallenge, "", ""},
		{"solo-mode-hello-world", &SeasonChallengeGet_Input{SeasonChallengeID: seasonChallenges["Solo Mode/Hello World"]}, nil, "Solo Mode", "Hello World"},
		{"no-team-in-season", &SeasonChallengeGet_Input{SeasonChallengeID: seasonChallenges["Test Season/Hello World"]}, errcode.ErrGetUserTeamFromSeason, "Test Season", "Hello World"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.SeasonChallengeGet(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			sc := ret.Item
			assert.Equal(t, test.input.SeasonChallengeID, sc.ID)
			assert.Equal(t, test.expectedChallengeName, sc.Flavor.Challenge.Name)
			assert.Equal(t, test.expectedSeasonName, sc.Season.Name)
		})
	}
}
