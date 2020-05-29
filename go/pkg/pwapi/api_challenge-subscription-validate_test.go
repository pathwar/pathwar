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
	challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{SeasonID: solo.ID})
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
		expectedValidations   int
	}{
		{"nil", nil, errcode.ErrMissingInput, "", 0},
		{"empty", &ChallengeSubscriptionValidate_Input{}, errcode.ErrMissingInput, "", 0},
		{"invalid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: 42, Passphrases: []string{"secret"}, Comment: "explanation"}, errcode.ErrGetChallengeSubscription, "", 0},
		{"valid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription.ChallengeSubscription.ID, Passphrases: []string{"secret"}, Comment: "ultra cool explanation"}, nil, "test", 1},
		// FIXME: revalidate
		// FIXME: new validation
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
		assert.Equalf(t, test.input.Passphrases, ret.ChallengeValidation.Passphrases, test.name)
		assert.Equalf(t, test.expectedPassphraseKey, ret.ChallengeValidation.PassphraseKey, test.name)
		assert.NotEmptyf(t, ret.ChallengeValidation.ChallengeSubscription.Validations, test.name)
		// fmt.Println(godev.PrettyJSON(ret))

		{
			chal, err := svc.SeasonChallengeGet(ctx, &SeasonChallengeGet_Input{SeasonChallengeID: challenges.Items[0].ID})
			require.NoErrorf(t, err, test.name)
			assert.Lenf(t, chal.Item.Subscriptions[0].Validations, test.expectedValidations, test.name)
			if len(chal.Item.Subscriptions[0].Validations) == 0 {
				continue
			}
			latest := chal.Item.Subscriptions[0].Validations[len(chal.Item.Subscriptions[0].Validations)-1]
			assert.Equalf(t, test.input.Comment, latest.AuthorComment, test.name)
			assert.Equalf(t, session.User.ID, latest.AuthorID, test.name)
			assert.Equalf(t, pwdb.ChallengeValidation_NeedReview, latest.Status, test.name)
			// FIXME: hide previous passphrases
			// fmt.Println(godev.PrettyJSON(chal))
		}
	}
}
