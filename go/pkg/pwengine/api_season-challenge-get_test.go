package pwengine

import (
	"context"
	"fmt"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
)

func TestEngine_SeasonChallengeGet(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := engine.UserGetSession(ctx, nil)
	checkErr(t, "", err)

	seasonChallenges := map[string]int64{}
	for _, seasonChallenge := range testingSeasonChallenges(t, engine).Items {
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
			ret, err := engine.SeasonChallengeGet(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			sc := ret.Item
			testSameInt64s(t, "", test.input.SeasonChallengeID, sc.ID)
			testSameStrings(t, "", test.expectedChallengeName, sc.Flavor.Challenge.Name)
			testSameStrings(t, "", test.expectedSeasonName, sc.Season.Name)
		})
	}
}
