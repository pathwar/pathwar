package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_ChallengeBuy(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	solo := testingSoloSeason(t, engine)

	// fetch challenges
	challenges, err := engine.SeasonChallengeList(ctx, &SeasonChallengeListInput{solo.ID})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// fetch user session
	session, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	activeTeam := session.User.ActiveTeamMember.Team

	var tests = []struct {
		name        string
		input       *SeasonChallengeBuyInput
		expectedErr error
	}{
		{"nil", nil, ErrMissingArgument},
		{"empty", &SeasonChallengeBuyInput{}, ErrMissingArgument},
		{
			"invalid season challenge ID",
			&SeasonChallengeBuyInput{SeasonChallengeID: 42, TeamID: activeTeam.ID},
			ErrInvalidArgument,
		},
		{
			"invalid team ID",
			&SeasonChallengeBuyInput{SeasonChallengeID: challenges.Items[0].ID, TeamID: 42},
			ErrInvalidArgument,
		},
		{
			"valid 1",
			&SeasonChallengeBuyInput{SeasonChallengeID: challenges.Items[0].ID, TeamID: activeTeam.ID},
			nil,
		},
		{
			"valid 2 (duplicate)",
			&SeasonChallengeBuyInput{SeasonChallengeID: challenges.Items[0].ID, TeamID: activeTeam.ID},
			ErrDuplicate,
		},
		// FIXME: check for a team and a challenge in different seasons
		// FIXME: check for a team from another user
		// FIXME: check for a challenge in draft mode
	}

	for _, test := range tests {
		subscription, err := engine.SeasonChallengeBuy(ctx, test.input)
		if !errors.Is(err, test.expectedErr) {
			t.Fatalf("%s: Expected %#v, got %#v.", test.name, test.expectedErr, err)
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
	}
}
