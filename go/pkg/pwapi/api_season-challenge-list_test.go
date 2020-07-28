package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestSvc_SeasonChallengeList(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	seasons := map[string]int64{}
	for _, season := range testingSeasons(t, svc).Items {
		seasons[season.Name] = season.ID
	}

	var tests = []struct {
		name          string
		input         *SeasonChallengeList_Input
		expectedErr   error
		expectedItems int
	}{
		{"empty", &SeasonChallengeList_Input{}, errcode.ErrMissingInput, 0},
		{"unknown-season-id", &SeasonChallengeList_Input{SeasonID: -42}, errcode.ErrInvalidSeasonID, 0},
		{"global-mode", &SeasonChallengeList_Input{SeasonID: seasons["Global"]}, nil, 7},
		{"test-season", &SeasonChallengeList_Input{SeasonID: seasons["Unit Test Season"]}, errcode.ErrUserHasNoTeamForSeason, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.SeasonChallengeList(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			//fmt.Println(godev.PrettyJSON(ret.Items))
			for _, item := range ret.Items {
				assert.Equal(t, test.input.SeasonID, item.SeasonID)
			}
			assert.Len(t, ret.Items, test.expectedItems)
		})
	}
}
