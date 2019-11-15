package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_SeasonChallengeList(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	seasons := map[string]int64{}
	for _, season := range testingSeasons(t, engine).Items {
		seasons[season.Name] = season.ID
	}

	var tests = []struct {
		name          string
		input         *SeasonChallengeListInput
		expectedErr   error
		expectedItems int
	}{
		{
			"empty",
			&SeasonChallengeListInput{},
			ErrMissingArgument,
			0,
		}, {
			"unknown-season-id",
			&SeasonChallengeListInput{SeasonID: -42}, // -42 should not exists
			ErrInvalidArgument,
			0,
		}, {
			"solo-mode",
			&SeasonChallengeListInput{SeasonID: seasons["Solo Mode"]},
			nil,
			5,
		}, {
			"test-season",
			&SeasonChallengeListInput{SeasonID: seasons["Test Season"]},
			ErrInvalidArgument,
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.SeasonChallengeList(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			//fmt.Println(godev.PrettyJSON(ret.Items))
			for _, item := range ret.Items {
				if item.SeasonID != test.input.SeasonID {
					t.Errorf("Expected %q, got %q.", test.input.SeasonID, item.SeasonID)
				}
			}
			if len(ret.Items) != test.expectedItems {
				t.Errorf("Expected %d, got %d.", test.expectedItems, len(ret.Items))
			}
		})
	}
}
