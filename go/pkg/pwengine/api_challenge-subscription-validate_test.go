package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/pwdb"
)

func TestEngine_ChallengeSubscriptionValidate(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	solo := testingSoloSeason(t, engine)

	// fetch user session
	session, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	activeTeam := session.User.ActiveTeamMember.Team

	// fetch challenges
	challenges, err := engine.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// buy a challenge
	subscription, err := engine.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	var tests = []struct {
		name                  string
		input                 *ChallengeSubscriptionValidate_Input
		expectedErr           error
		expectedPassphraseKey string
	}{
		{"nil", nil, ErrMissingArgument, ""},
		{"empty", &ChallengeSubscriptionValidate_Input{}, ErrMissingArgument, ""},
		{"invalid", &ChallengeSubscriptionValidate_Input{ChallengeSubscriptionID: 42, Passphrase: "secret", Comment: "explanation"}, ErrInvalidArgument, ""},
		{
			"valid",
			&ChallengeSubscriptionValidate_Input{
				ChallengeSubscriptionID: subscription.ChallengeSubscription.ID,
				Passphrase:              "secret",
				Comment:                 "ultra cool explanation",
			},
			nil,
			"test",
		},
	}

	for _, test := range tests {
		ret, err := engine.ChallengeSubscriptionValidate(ctx, test.input)
		if test.expectedErr != err {
			t.Errorf("%s: Expected %v, got %v.", test.name, test.expectedErr, err)
		}
		if err != nil {
			continue
		}

		if ret.ChallengeValidation.ChallengeSubscriptionID != subscription.ChallengeSubscription.ID {
			t.Errorf("%s: Expected %d, got %d.", test.name, subscription.ChallengeSubscription.ID, ret.ChallengeValidation.ChallengeSubscriptionID)
		}
		if ret.ChallengeValidation.AuthorID != session.User.ID {
			t.Errorf("%s: Expectd %d, got %d.", test.name, session.User.ID, ret.ChallengeValidation.AuthorID)
		}
		if ret.ChallengeValidation.Status != pwdb.ChallengeValidation_NeedReview {
			t.Errorf("%s: Expected %v, got %v.", test.name, pwdb.ChallengeValidation_NeedReview, ret.ChallengeValidation.Status)
		}
		if test.input.Comment != ret.ChallengeValidation.AuthorComment {
			t.Errorf("%s: Expected %q, got %q.", test.name, test.input.Comment, ret.ChallengeValidation.AuthorComment)
		}
		if test.input.Passphrase != ret.ChallengeValidation.Passphrase {
			t.Errorf("%s: Expected %q, got %q.", test.name, test.input.Passphrase, ret.ChallengeValidation.Passphrase)
		}
		if test.expectedPassphraseKey != ret.ChallengeValidation.PassphraseKey {
			t.Errorf("%s: Expected %v, got %v.", test.name, test.expectedPassphraseKey, ret.ChallengeValidation.PassphraseKey)
		}
		if len(ret.ChallengeValidation.ChallengeSubscription.Validations) == 0 {
			t.Errorf("%s: should have at least one validation", test.name)
		}
		// fmt.Println(godev.PrettyJSON(ret))
	}
}
