package pwapi

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
)

func TestService_ChallengeGet(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	challenges := map[string]int64{}
	for _, challenge := range testingChallenges(t, svc).Items {
		challenges[challenge.Name] = challenge.ID
	}

	var tests = []struct {
		name                  string
		input                 *ChallengeGet_Input
		expectedErr           error
		expectedChallengeName string
		expectedAuthor        string
	}{
		{"empty", &ChallengeGet_Input{}, errcode.ErrMissingInput, "", ""},
		{"unknown-season-id", &ChallengeGet_Input{ChallengeID: -42}, errcode.ErrDBNotFound, "", ""}, // -42 should not exists
		{"Staff", &ChallengeGet_Input{ChallengeID: challenges["Hello World"]}, nil, "Hello World", "Staff Team"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.ChallengeGet(ctx, test.input)
			testSameErrcodes(t, test.name, test.expectedErr, err)
			if err != nil {
				return
			}

			// FIXME: check for ChallengeVersions and ChallengeInstances
			testSameInt64s(t, test.name, test.input.ChallengeID, ret.Item.ID)
			testSameStrings(t, test.name, test.expectedChallengeName, ret.Item.Name)
			testSameStrings(t, test.name, test.expectedAuthor, ret.Item.Author)
		})
	}
}
