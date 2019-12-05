package pwapi

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
)

func TestSvc_UserSetPreferences(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// get user session before setting preferences
	beforeSession, err := svc.UserGetSession(ctx, nil)
	checkErr(t, "", err)
	seasons := map[string]int64{}
	for _, season := range beforeSession.Seasons {
		seasons[season.Season.Name] = season.Season.ID
	}

	var tests = []struct {
		name                 string
		input                *UserSetPreferences_Input
		expectedErr          error
		expectedSeasonID     int64
		expectedTeamMemberID int64
	}{
		{"empty", &UserSetPreferences_Input{}, errcode.ErrMissingInput, seasons["Solo Mode"], beforeSession.User.ActiveTeamMemberID},
		{"unknown-season-id", &UserSetPreferences_Input{ActiveSeasonID: -42}, errcode.ErrInvalidSeasonID, seasons["Solo Mode"], beforeSession.User.ActiveTeamMemberID},
		{"solo-mode", &UserSetPreferences_Input{ActiveSeasonID: seasons["Solo Mode"]}, nil, seasons["Solo Mode"], beforeSession.User.ActiveTeamMemberID},
		{"test-season", &UserSetPreferences_Input{ActiveSeasonID: seasons["Test Season"]}, nil, seasons["Test Season"], 0},
		{"solo-mode-again", &UserSetPreferences_Input{ActiveSeasonID: seasons["Solo Mode"]}, nil, seasons["Solo Mode"], beforeSession.User.ActiveTeamMemberID},
	}

	for _, test := range tests {
		_, err := svc.UserSetPreferences(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)

		session, err := svc.UserGetSession(ctx, nil)
		checkErr(t, test.name, err)

		testSameInt64s(t, test.name, test.expectedSeasonID, session.User.ActiveSeasonID)
		testSameInt64s(t, test.name, test.expectedTeamMemberID, session.User.ActiveTeamMemberID)
	}
}
