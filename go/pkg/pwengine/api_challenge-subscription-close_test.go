package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func TestEngine_ChallengeSubscriptionClose(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	solo := testingSoloSeason(t, engine)

	// fetch user session
	session, err := engine.UserGetSession(ctx, nil)
	checkErr(t, "", err)
	activeTeam := session.User.ActiveTeamMember.Team

	// fetch challenges
	challenges, err := engine.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
	checkErr(t, "", err)

	// buy two challenges
	subscription1, err := engine.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	checkErr(t, "", err)
	subscription2, err := engine.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[1].ID,
		TeamID:            activeTeam.ID,
	})
	checkErr(t, "", err)

	// validate second challenge
	_, err = engine.ChallengeSubscriptionValidate(ctx, &ChallengeSubscriptionValidate_Input{
		ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID,
		Passphrase:              "secret",
	})
	checkErr(t, "", err)

	var tests = []struct {
		name        string
		input       *ChallengeSubscriptionClose_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &ChallengeSubscriptionClose_Input{}, errcode.ErrMissingInput},
		{"subscription1", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription1.ChallengeSubscription.ID}, errcode.ErrMissingChallengeValidation},
		{"subscription2", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID}, nil},
		{"subscription2-again", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID}, errcode.ErrChallengeAlreadyClosed},
	}
	for _, test := range tests {
		ret, err := engine.ChallengeSubscriptionClose(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)
		if err != nil {
			continue
		}

		testIsNotNil(t, test.name, ret.ChallengeSubscription.ClosedAt)
		testSameInt64s(t, test.name, session.User.ID, ret.ChallengeSubscription.CloserID)
		testSameAnys(t, test.name, pwdb.ChallengeSubscription_Closed, ret.ChallengeSubscription.Status)
		testSameInt64s(t, test.name, activeTeam.ID, ret.ChallengeSubscription.Team.ID)
		testSameInt64s(t, test.name, test.input.ChallengeSubscriptionID, ret.ChallengeSubscription.ID)
		if len(ret.ChallengeSubscription.Validations) == 0 {
			t.Errorf("%s: should have at least one validation", test.name)
		}
	}
}
