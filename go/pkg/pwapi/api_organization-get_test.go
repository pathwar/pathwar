package pwapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestService_OrganizationGet(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	organizations := map[string]int64{}
	for _, organization := range testingOrganizations(t, svc).Items {
		key := fmt.Sprintf("%s", organization.Name)
		organizations[key] = organization.ID
	}

	tests := []struct {
		name                     string
		input                    *OrganizationGet_Input
		expectedErr              error
		expectedOrganizationName string
		expectedSeasonName       string
	}{
		{"empty", &OrganizationGet_Input{}, errcode.ErrMissingInput, "", ""},
		{"unknown-season-id", &OrganizationGet_Input{OrganizationID: -42}, errcode.ErrGetOrganization, "", ""},
		{"Staff", &OrganizationGet_Input{OrganizationID: organizations["Staff"]}, nil, "Staff", "Global"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.OrganizationGet(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, test.input.OrganizationID, ret.Item.ID)
			assert.Equal(t, test.expectedOrganizationName, ret.Item.Name)

			seasonName := ""
			for _, team := range ret.Item.Teams {
				if team.Season.Name == test.expectedSeasonName {
					seasonName = team.Season.Name
				}
			}
			assert.Equal(t, test.expectedSeasonName, seasonName)
		})
	}
}
