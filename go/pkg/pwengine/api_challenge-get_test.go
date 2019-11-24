package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_ChallengeGet(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	challenges := map[string]int64{}
	for _, challenge := range testingChallenges(t, engine).Items {
		challenges[challenge.Name] = challenge.ID
	}

	var tests = []struct {
		name                  string
		input                 *ChallengeGet_Input
		expectedErr           error
		expectedChallengeName string
		expectedAuthor        string
	}{
		{
			"empty",
			&ChallengeGet_Input{},
			ErrMissingArgument,
			"",
			"",
		}, {
			"unknown-season-id",
			&ChallengeGet_Input{ChallengeID: -42}, // -42 should not exists
			ErrInvalidArgument,
			"",
			"",
		}, {
			"Staff",
			&ChallengeGet_Input{ChallengeID: challenges["Hello World"]},
			nil,
			"Hello World",
			"Staff Team",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.ChallengeGet(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			// FIXME: check for ChallengeVersions and ChallengeInstances

			if ret.Item.ID != test.input.ChallengeID {
				t.Fatalf("Expected %q, got %q.", test.input.ChallengeID, ret.Item.ID)
			}
			if ret.Item.Name != test.expectedChallengeName {
				t.Fatalf("Expected %q, got %q.", test.expectedChallengeName, ret.Item.Name)
			}
			if ret.Item.Author != test.expectedAuthor {
				t.Fatalf("Expected %q, got %q.", test.expectedAuthor, ret.Item.Author)
			}
		})
	}
}
