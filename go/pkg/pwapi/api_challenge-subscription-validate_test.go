package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
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

	// buy challenges
	subscription1, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	require.NoError(t, err)
	subscription2, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[9].ID,
		TeamID:            activeTeam.ID,
	})
	require.NoError(t, err)

	var tests = []struct {
		name                            string
		input                           *ChallengeSubscriptionValidate_Input
		expectedErr                     error
		expectedPassphraseIndices       string
		expectedChallengeSubscriptionID int64
		expectedValidations             int
		expectedCash                    int64
	}{
		{"nil", nil, errcode.ErrMissingInput, "", 0, 0, 0},
		{"empty", &ChallengeSubscriptionValidate_Input{}, errcode.ErrMissingInput, "", 0, 0, 0},
		{"invalid challenge subscription", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: 42, Passphrases: []string{"secret"}, Comment: "explanation"}, errcode.ErrGetChallengeSubscription, "", 0, 0, 0},
		{"challenge with no instance", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription1.ChallengeSubscription.ID, Passphrases: []string{"secret"}, Comment: "ultra cool explanation"}, errcode.ErrChallengeInactiveValidation, "test", 0, 0, 0},
		{"3/4 valid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID, Passphrases: []string{"a", "b", "c"}, Comment: "ultra cool explanation"}, errcode.ErrChallengeIncompleteValidation, "", 0, 0, 0},
		{"too many", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID, Passphrases: []string{"a", "b", "c", "d", "e"}, Comment: "ultra cool explanation"}, errcode.ErrChallengeIncompleteValidation, "", 0, 0, 0},
		{"bad passphrases", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID, Passphrases: []string{"e", "f", "g", "h"}, Comment: "ultra cool explanation"}, errcode.ErrChallengeIncompleteValidation, "", 0, 0, 0},
		{"one good passphrase", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID, Passphrases: []string{"a", "f", "g", "h"}, Comment: "ultra cool explanation"}, errcode.ErrChallengeIncompleteValidation, "", 0, 0, 0},
		{"valid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID, Passphrases: []string{"a", "b", "c", "d"}, Comment: "ultra cool explanation"}, nil, "[0,1,2,3]", subscription2.ChallengeSubscription.ID, 1, 10},
		// {"revalid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID, Passphrases: []string{"a", "b", "c", "d"}, Comment: "ultra cool explanation"}, nil, "test", 1},
		// FIXME: new validation
	}

	for _, test := range tests {
		ret, err := svc.ChallengeSubscriptionValidate(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.expectedChallengeSubscriptionID, ret.ChallengeValidation.ChallengeSubscriptionID, test.name)
		assert.Equalf(t, session.User.ID, ret.ChallengeValidation.AuthorID, test.name)
		assert.Equalf(t, pwdb.ChallengeValidation_NeedReview, ret.ChallengeValidation.Status, test.name)
		assert.Equalf(t, test.input.Comment, ret.ChallengeValidation.AuthorComment, test.name)
		assert.Equalf(t, test.expectedPassphraseIndices, ret.ChallengeValidation.Passphrases, test.name)
		assert.NotEmptyf(t, ret.ChallengeValidation.ChallengeSubscription.Validations, test.name)
		assert.NotNilf(t, ret.ChallengeValidation.ChallengeSubscription.ClosedAt, test.name)
		assert.Equalf(t, session.User.ID, ret.ChallengeValidation.ChallengeSubscription.CloserID, test.name)
		assert.Equalf(t, pwdb.ChallengeSubscription_Closed, ret.ChallengeValidation.ChallengeSubscription.Status, test.name)
		assert.Equalf(t, activeTeam.ID, ret.ChallengeValidation.ChallengeSubscription.Team.ID, test.name)
		assert.Equalf(t, test.input.ChallengeSubscriptionID, ret.ChallengeValidation.ChallengeSubscription.ID, test.name)
		assert.NotEmptyf(t, ret.ChallengeValidation.ChallengeSubscription.Validations, test.name)
		assert.Equal(t, ret.ChallengeValidation.ChallengeSubscription.Team.Cash, test.expectedCash, test.name)

		{
			challenge, err := svc.SeasonChallengeGet(ctx, &SeasonChallengeGet_Input{SeasonChallengeID: ret.ChallengeValidation.ChallengeSubscription.SeasonChallenge.ID})
			require.NoErrorf(t, err, test.name)
			assert.Lenf(t, challenge.Item.Subscriptions[0].Validations, test.expectedValidations, test.name)
			if len(challenge.Item.Subscriptions[0].Validations) == 0 {
				continue
			}
			latest := challenge.Item.Subscriptions[0].Validations[len(challenge.Item.Subscriptions[0].Validations)-1]
			assert.Equalf(t, test.input.Comment, latest.AuthorComment, test.name)
			assert.Equalf(t, session.User.ID, latest.AuthorID, test.name)
			assert.Equalf(t, pwdb.ChallengeValidation_NeedReview, latest.Status, test.name)
			// FIXME: hide previous passphrases
		}
	}
}
