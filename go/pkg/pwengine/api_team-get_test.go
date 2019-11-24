package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_TeamGet(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	organizations := map[string]int64{}
	for _, organization := range testingTeams(t, engine).Items {
		organizations[organization.Organization.Name] = organization.ID
	}

	var tests = []struct {
		name                     string
		input                    *TeamGet_Input
		expectedErr              error
		expectedOrganizationName string
		expectedSeasonName       string
	}{
		{
			"empty",
			&TeamGet_Input{},
			ErrMissingArgument,
			"",
			"",
		}, {
			"unknown-season-id",
			&TeamGet_Input{TeamID: -42}, // -42 should not exists
			ErrInvalidArgument,
			"",
			"",
		}, {
			"Staff",
			&TeamGet_Input{TeamID: organizations["Staff"]},
			nil,
			"Staff",
			"Solo Mode",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.TeamGet(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			if ret.Item.ID != test.input.TeamID {
				t.Fatalf("Expected %q, got %q.", test.input.TeamID, ret.Item.ID)
			}
			if ret.Item.Organization.Name != test.expectedOrganizationName {
				t.Fatalf("Expected %q, got %q.", test.expectedOrganizationName, ret.Item.Organization.Name)
			}
			if ret.Item.Season.Name != test.expectedSeasonName {
				t.Fatalf("Expected %q, got %q.", test.expectedSeasonName, ret.Item.Season.Name)
			}
		})
	}
}
