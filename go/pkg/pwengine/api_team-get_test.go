package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
)

func TestEngine_TeamGet(t *testing.T) {
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
		{"empty", &TeamGet_Input{}, errcode.ErrMissingInput, "", ""},
		{"unknown-season-id", &TeamGet_Input{TeamID: -42}, errcode.ErrGetTeam, "", ""},
		{"Staff", &TeamGet_Input{TeamID: organizations["Staff"]}, nil, "Staff", "Solo Mode"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.TeamGet(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			testSameInt64s(t, "", test.input.TeamID, ret.Item.ID)
			testSameStrings(t, "", test.expectedOrganizationName, ret.Item.Organization.Name)
			testSameStrings(t, "", test.expectedSeasonName, ret.Item.Season.Name)
		})
	}
}
