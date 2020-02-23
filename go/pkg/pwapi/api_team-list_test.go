package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/v2/go/internal/testutil"
	"pathwar.land/v2/go/pkg/errcode"
)

func TestSvc_TeamList(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	seasons := map[string]int64{}
	for _, season := range testingSeasons(t, svc).Items {
		seasons[season.Name] = season.ID
	}

	var tests = []struct {
		name                  string
		input                 *TeamList_Input
		expectedErr           error
		expectedOrganizations int
		// expectedOwnedOrganizations int?
	}{
		{"empty", &TeamList_Input{}, errcode.ErrMissingInput, 0},
		{"unknown-season-id", &TeamList_Input{SeasonID: -42}, errcode.ErrInvalidSeasonID, 0},
		{"solo-mode", &TeamList_Input{SeasonID: seasons["Solo Mode"]}, nil, 1},
		{"test-season", &TeamList_Input{SeasonID: seasons["Test Season"]}, nil, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.TeamList(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, test.expectedOrganizations, len(ret.Items))

			// fmt.Println(godev.PrettyJSON(ret))
			for _, organization := range ret.Items {
				assert.Equal(t, test.input.SeasonID, organization.SeasonID)
			}
		})
	}
}
