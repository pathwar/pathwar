package pwapi

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func TestSvc_ChallengeSubscriptionValidate(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	solo := testingSoloSeason(t, svc)

	// fetch user session
	session, err := svc.UserGetSession(ctx, nil)
	checkErr(t, "", err)
	activeTeam := session.User.ActiveTeamMember.Team

	// fetch challenges
	challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
	checkErr(t, "", err)

	// buy a challenge
	subscription, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	checkErr(t, "", err)

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

		testSameInt64s(t, test.name, subscription.ChallengeSubscription.ID, ret.ChallengeValidation.ChallengeSubscriptionID)
		testSameInt64s(t, test.name, session.User.ID, ret.ChallengeValidation.AuthorID)
		testSameAnys(t, test.name, pwdb.ChallengeValidation_NeedReview, ret.ChallengeValidation.Status)
		testSameStrings(t, test.name, test.input.Comment, ret.ChallengeValidation.AuthorComment)
		testSameStrings(t, test.name, test.input.Passphrase, ret.ChallengeValidation.Passphrase)
		testSameStrings(t, test.name, test.expectedPassphraseKey, ret.ChallengeValidation.PassphraseKey)
		if len(ret.ChallengeValidation.ChallengeSubscription.Validations) == 0 {
			t.Errorf("%s: should have at least one validation", test.name)
		}
		// fmt.Println(godev.PrettyJSON(ret))
	}
}
