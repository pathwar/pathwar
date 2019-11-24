package pwengine

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_SeasonChallengeGet(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

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
		{
			"empty",
			&SeasonChallengeGet_Input{},
			ErrMissingArgument,
			"",
			"",
		}, {
			"unknown-season-id",
			&SeasonChallengeGet_Input{SeasonChallengeID: -42}, // -42 should not exists
			ErrInvalidArgument,
			"",
			"",
		}, {
			"solo-mode-hello-world",
			&SeasonChallengeGet_Input{SeasonChallengeID: seasonChallenges["Solo Mode/Hello World"]},
			nil,
			"Solo Mode",
			"Hello World",
		}, {
			"test-season-hello-world",
			&SeasonChallengeGet_Input{SeasonChallengeID: seasonChallenges["Test Season/Hello World"]},
			ErrInvalidArgument, // no team in this season
			"Test Season",
			"Hello World",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.SeasonChallengeGet(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			//fmt.Println(godev.PrettyJSON(ret.Item))
			sc := ret.Item
			if test.input.SeasonChallengeID != sc.ID {
				t.Errorf("Expected %d, got %d.", test.input.SeasonChallengeID, sc.ID)
			}
			if sc.Flavor.Challenge.Name != test.expectedChallengeName {
				t.Errorf("Expected %q, got %q.", test.expectedChallengeName, sc.Flavor.Challenge.Name)
			}
			if sc.Season.Name != test.expectedSeasonName {
				t.Errorf("Expected %q, got %q.", test.expectedSeasonName, sc.Season.Name)
			}
		})
	}
}
