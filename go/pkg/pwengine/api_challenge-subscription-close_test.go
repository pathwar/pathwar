package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/pwdb"
)

func TestEngine_ChallengeSubscriptionClose(t *testing.T) {
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

	// buy two challenges
	subscription1, err := engine.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	subscription2, err := engine.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[1].ID,
		TeamID:            activeTeam.ID,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// validate second challenge
	_, err = engine.ChallengeSubscriptionValidate(ctx, &ChallengeSubscriptionValidate_Input{
		ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID,
		Passphrase:              "secret",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	var tests = []struct {
		name        string
		input       *ChallengeSubscriptionClose_Input
		expectedErr error
	}{
		{"nil", nil, ErrMissingArgument},
		{"empty", &ChallengeSubscriptionClose_Input{}, ErrMissingArgument},
		{"subscription1", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription1.ChallengeSubscription.ID}, ErrMissingRequiredValidation},
		{"subscription2", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID}, nil},
		{"subscription2", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID}, ErrInvalidArgument},
	}
	for _, test := range tests {
		ret, err := engine.ChallengeSubscriptionClose(ctx, test.input)
		if test.expectedErr != err {
			t.Errorf("%s: Expected %v, got %v.", test.name, test.expectedErr, err)
		}
		if err != nil {
			continue
		}

		if ret.ChallengeSubscription.ClosedAt == nil {
			t.Errorf("%s: Expected ClosedAt != nil.", test.name)
		}
		if ret.ChallengeSubscription.CloserID != session.User.ID {
			t.Errorf("%s: Expected %d, got %d.", test.name, session.User.ID, ret.ChallengeSubscription.CloserID)
		}
		if ret.ChallengeSubscription.Status != pwdb.ChallengeSubscription_Closed {
			t.Errorf("%s: Expected %v, got %v.", test.name, pwdb.ChallengeSubscription_Closed, ret.ChallengeSubscription.Status)
		}
		if ret.ChallengeSubscription.Team.ID != activeTeam.ID {
			t.Errorf("%s: Expected %d, got %d.", test.name, activeTeam.ID, ret.ChallengeSubscription.Team.ID)
		}
		if test.input.ChallengeSubscriptionID != ret.ChallengeSubscription.ID {
			t.Errorf("%s: Expected %d, got %d.", test.name, test.input.ChallengeSubscriptionID, ret.ChallengeSubscription.ID)
		}
		if len(ret.ChallengeSubscription.Validations) == 0 {
			t.Errorf("%s: should have at least one validation", test.name)
		}
	}
}
