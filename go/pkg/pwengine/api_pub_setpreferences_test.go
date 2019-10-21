package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_SetPreferences(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// get user session before setting preferences
	beforeSession, err := engine.GetUserSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tournaments := map[string]int64{}
	for _, tournament := range beforeSession.Tournaments {
		tournaments[tournament.Tournament.Name] = tournament.Tournament.ID
	}

	var tests = []struct {
		name                       string
		input                      *SetPreferencesInput
		expectedErr                error
		expectedTournamentID       int64
		expectedTournamentMemberID int64
	}{
		{
			"empty",
			&SetPreferencesInput{},
			ErrMissingArgument,
			tournaments["Solo Mode"],
			beforeSession.User.ActiveTournamentMemberID,
		}, {
			"unknown-tournament-id",
			&SetPreferencesInput{ActiveTournamentID: -42}, // should not exists
			ErrInvalidArgument,
			tournaments["Solo Mode"],
			beforeSession.User.ActiveTournamentMemberID,
		}, {
			"solo-mode",
			&SetPreferencesInput{ActiveTournamentID: tournaments["Solo Mode"]},
			nil,
			tournaments["Solo Mode"],
			beforeSession.User.ActiveTournamentMemberID,
		}, {
			"test-tournament",
			&SetPreferencesInput{ActiveTournamentID: tournaments["Test Tournament"]},
			nil,
			tournaments["Test Tournament"],
			0,
		}, {
			"solo-mode-again",
			&SetPreferencesInput{ActiveTournamentID: tournaments["Solo Mode"]},
			nil,
			tournaments["Solo Mode"],
			beforeSession.User.ActiveTournamentMemberID,
		},
	}

	for _, test := range tests {
		_, err := engine.SetPreferences(ctx, test.input)
		if !errors.Is(err, test.expectedErr) {
			t.Fatalf("%s: Expected %#v, got %#v.", test.name, test.expectedErr, err)
		}

		session, err := engine.GetUserSession(ctx, nil)
		if err != nil {
			t.Fatalf("%s: err: %v", test.name, err)
		}

		if session.User.ActiveTournamentID != test.expectedTournamentID {
			t.Fatalf("%s: Expected %d, got %d.", test.name, test.expectedTournamentID, session.User.ActiveTournamentID)
		}
		if session.User.ActiveTournamentMemberID != test.expectedTournamentMemberID {
			t.Fatalf("%s: Expected %d, got %d.", test.name, test.expectedTournamentMemberID, session.User.ActiveTournamentMemberID)
		}
	}
}
