package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_UserSetPreferences(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// get user session before setting preferences
	beforeSession, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	seasons := map[string]int64{}
	for _, season := range beforeSession.Seasons {
		seasons[season.Season.Name] = season.Season.ID
	}

	var tests = []struct {
		name                 string
		input                *UserSetPreferencesInput
		expectedErr          error
		expectedSeasonID     int64
		expectedTeamMemberID int64
	}{
		{
			"empty",
			&UserSetPreferencesInput{},
			ErrMissingArgument,
			seasons["Solo Mode"],
			beforeSession.User.ActiveTeamMemberID,
		}, {
			"unknown-season-id",
			&UserSetPreferencesInput{ActiveSeasonID: -42}, // should not exists
			ErrInvalidArgument,
			seasons["Solo Mode"],
			beforeSession.User.ActiveTeamMemberID,
		}, {
			"solo-mode",
			&UserSetPreferencesInput{ActiveSeasonID: seasons["Solo Mode"]},
			nil,
			seasons["Solo Mode"],
			beforeSession.User.ActiveTeamMemberID,
		}, {
			"test-season",
			&UserSetPreferencesInput{ActiveSeasonID: seasons["Test Season"]},
			nil,
			seasons["Test Season"],
			0,
		}, {
			"solo-mode-again",
			&UserSetPreferencesInput{ActiveSeasonID: seasons["Solo Mode"]},
			nil,
			seasons["Solo Mode"],
			beforeSession.User.ActiveTeamMemberID,
		},
	}

	for _, test := range tests {
		_, err := engine.UserSetPreferences(ctx, test.input)
		if !errors.Is(err, test.expectedErr) {
			t.Fatalf("%s: Expected %#v, got %#v.", test.name, test.expectedErr, err)
		}

		session, err := engine.UserGetSession(ctx, nil)
		if err != nil {
			t.Fatalf("%s: err: %v", test.name, err)
		}

		if session.User.ActiveSeasonID != test.expectedSeasonID {
			t.Fatalf("%s: Expected %d, got %d.", test.name, test.expectedSeasonID, session.User.ActiveSeasonID)
		}
		if session.User.ActiveTeamMemberID != test.expectedTeamMemberID {
			t.Fatalf("%s: Expected %d, got %d.", test.name, test.expectedTeamMemberID, session.User.ActiveTeamMemberID)
		}
	}
}
