package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_ChallengeBuy(t *testing.T) {
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

	var tests = []struct {
		name        string
		input       *SeasonChallengeBuy_Input
		expectedErr error
	}{
		{"nil", nil, ErrMissingArgument},
		{"empty", &SeasonChallengeBuy_Input{}, ErrMissingArgument},
		{
			"invalid season challenge ID",
			&SeasonChallengeBuy_Input{SeasonChallengeID: 42, TeamID: activeTeam.ID},
			ErrInvalidArgument,
		},
		{
			"invalid team ID",
			&SeasonChallengeBuy_Input{SeasonChallengeID: challenges.Items[0].ID, TeamID: 42},
			ErrInvalidArgument,
		},
		{
			"valid 1",
			&SeasonChallengeBuy_Input{SeasonChallengeID: challenges.Items[0].ID, TeamID: activeTeam.ID},
			nil,
		},
		{
			"valid 2 (duplicate)",
			&SeasonChallengeBuy_Input{SeasonChallengeID: challenges.Items[0].ID, TeamID: activeTeam.ID},
			ErrDuplicate,
		},
		// FIXME: check for a team and a challenge in different seasons
		// FIXME: check for a team from another user
		// FIXME: check for a challenge in draft mode
	}

	for _, test := range tests {
		subscription, err := engine.SeasonChallengeBuy(ctx, test.input)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: Expected %#v, got %#v.", test.name, test.expectedErr, err)
		}

		if err != nil {
			continue
		}
		if subscription.ChallengeSubscription.TeamID != test.input.TeamID {
			t.Errorf("%s: Expected %d, got %d.", test.name, test.input.TeamID, subscription.ChallengeSubscription.TeamID)
		}
		if subscription.ChallengeSubscription.SeasonChallengeID != test.input.SeasonChallengeID {
			t.Errorf("%s: Expected %d, got %d.", test.name, test.input.SeasonChallengeID, subscription.ChallengeSubscription.SeasonChallengeID)
		}
		if subscription.ChallengeSubscription.BuyerID != session.User.ID {
			t.Errorf("%s: Expected %d, got %d.", test.name, session.User.ID, subscription.ChallengeSubscription.BuyerID)
		}

		// check if challenge subscription is now visible in season challenge list
		challenges, err := engine.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		found := 0
		for _, challenge := range challenges.Items {
			if challenge.ID == subscription.ChallengeSubscription.SeasonChallengeID {
				found++
				if len(challenge.Subscriptions) != 1 {
					t.Errorf("%s: Expected only one subscription, got %d.", test.name, len(challenge.Subscriptions))
				}

				if challenge.Subscriptions[0].ID != subscription.ChallengeSubscription.ID {
					t.Errorf("%s: Expected %d, got %d.", test.name, subscription.ChallengeSubscription.ID, challenge.Subscriptions[0].ID)
				}

				if challenge.Subscriptions[0].TeamID != test.input.TeamID {
					t.Errorf("%s: Expected %d, got %d.", test.name, test.input.TeamID, challenge.Subscriptions[0].TeamID)
				}
			}
		}
		if found != 1 {
			t.Errorf("%s: Expected 1 found, got %d.", test.name, found)
		}
	}
}
