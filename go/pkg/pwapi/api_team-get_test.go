package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestSvc_TeamGet(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	organizations := map[string]int64{}
	for _, organization := range testingTeams(t, svc).Items {
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
			ret, err := svc.TeamGet(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, test.input.TeamID, ret.Item.ID)
			assert.Equal(t, test.expectedOrganizationName, ret.Item.Organization.Name)
			assert.Equal(t, test.expectedSeasonName, ret.Item.Season.Name)
		})
	}
}
