package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_TeamList(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	seasons := map[string]int64{}
	for _, season := range testingSeasons(t, engine).Items {
		seasons[season.Name] = season.ID
	}

	var tests = []struct {
		name                  string
		input                 *TeamList_Input
		expectedErr           error
		expectedOrganizations int
		// expectedOwnedOrganizations int?
	}{
		{
			"empty",
			&TeamList_Input{},
			ErrMissingArgument,
			0,
		}, {
			"unknown-season-id",
			&TeamList_Input{SeasonID: -42}, // -42 should not exists
			ErrInvalidArgument,
			0,
		}, {
			"solo-mode",
			&TeamList_Input{SeasonID: seasons["Solo Mode"]},
			nil,
			1,
		}, {
			"test-season",
			&TeamList_Input{SeasonID: seasons["Test Season"]},
			nil,
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.TeamList(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			// fmt.Println(godev.PrettyJSON(ret))
			for _, organization := range ret.Items {
				if organization.SeasonID != test.input.SeasonID {
					t.Fatalf("Expected %q, got %q.", test.input.SeasonID, organization.SeasonID)
				}
			}

			if len(ret.Items) != test.expectedOrganizations {
				t.Fatalf("Expected %d, got %d.", test.expectedOrganizations, len(ret.Items))
			}
		})
	}
}
