package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/v2/go/internal/testutil"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func TestSvc_ChallengeSubscriptionValidate(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	solo := testingSoloSeason(t, svc)

	// fetch user session
	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
	activeTeam := session.User.ActiveTeamMember.Team

	// fetch challenges
	challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
	require.NoError(t, err)

	// buy a challenge
	subscription, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	require.NoError(t, err)

	var tests = []struct {
		name                  string
		input                 *ChallengeSubscriptionValidate_Input
		expectedErr           error
		expectedPassphraseKey string
	}{
		{"nil", nil, errcode.ErrMissingInput, ""},
		{"empty", &ChallengeSubscriptionValidate_Input{}, errcode.ErrMissingInput, ""},
		{"invalid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: 42, Passphrase: "secret", Comment: "explanation"}, errcode.ErrGetChallengeSubscription, ""},
		{"valid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription.ChallengeSubscription.ID, Passphrase: "secret", Comment: "ultra cool explanation"}, nil, "test"},
	}

	for _, test := range tests {
		ret, err := svc.ChallengeSubscriptionValidate(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)
		if err != nil {
			continue
		}

		assert.Equalf(t, subscription.ChallengeSubscription.ID, ret.ChallengeValidation.ChallengeSubscriptionID, test.name)
		assert.Equalf(t, session.User.ID, ret.ChallengeValidation.AuthorID, test.name)
		assert.Equalf(t, pwdb.ChallengeValidation_NeedReview, ret.ChallengeValidation.Status, test.name)
		assert.Equalf(t, test.input.Comment, ret.ChallengeValidation.AuthorComment, test.name)
		assert.Equalf(t, test.input.Passphrase, ret.ChallengeValidation.Passphrase, test.name)
		assert.Equalf(t, test.expectedPassphraseKey, ret.ChallengeValidation.PassphraseKey, test.name)
		assert.NotEmptyf(t, ret.ChallengeValidation.ChallengeSubscription.Validations, test.name)
		// fmt.Println(godev.PrettyJSON(ret))
	}
}
