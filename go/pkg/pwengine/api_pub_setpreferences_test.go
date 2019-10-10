package pwengine

import (
	"context"
	"errors"
	"testing"
)

func TestEngine_SetPreferences(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{})
	defer cleanup()
	ctx := testSetContextToken(t, context.Background())

	// get user session before setting preferences
	beforeSession, err := engine.GetUserSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tournaments := map[string]string{}
	for _, tournament := range beforeSession.Tournaments {
		tournaments[tournament.Tournament.Name] = tournament.Tournament.ID
	}

	var tests = []struct {
		name                       string
		input                      *SetPreferencesInput
		expectedErr                error
		expectedTournamentID       string
		expectedTournamentMemberID string
	}{
		{
			"empty",
			&SetPreferencesInput{},
			ErrMissingArgument,
			tournaments["Solo Mode"],
			beforeSession.User.ActiveTournamentMemberID,
		}, {
			"unknown-tournament-id",
			&SetPreferencesInput{ActiveTournamentID: "does not exist"},
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
			"",
		}, {
			"solo-mode",
			&SetPreferencesInput{ActiveTournamentID: tournaments["Solo Mode"]},
			nil,
			tournaments["Solo Mode"],
			beforeSession.User.ActiveTournamentMemberID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := engine.SetPreferences(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}

			session, err := engine.GetUserSession(ctx, nil)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			if session.User.ActiveTournamentID != test.expectedTournamentID {
				t.Fatalf("Expected %s, got %s.", test.expectedTournamentID, session.User.ActiveTournamentID)
			}
			if session.User.ActiveTournamentMemberID != test.expectedTournamentMemberID {
				t.Fatalf("Expected %s, got %s.", test.expectedTournamentMemberID, session.User.ActiveTournamentMemberID)
			}
		})
	}
}
