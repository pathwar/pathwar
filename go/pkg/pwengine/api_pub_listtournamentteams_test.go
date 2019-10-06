package pwengine

import (
	"context"
	"errors"
	"testing"
)

func TestEngine_ListTournamentTeams(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{})
	defer cleanup()
	ctx := testSetContextToken(t, context.Background())

	// FIXME: check for permissions

	tournaments := map[string]string{}
	for _, tournament := range testingTournaments(t, engine).Items {
		tournaments[tournament.Name] = tournament.ID
	}

	var tests = []struct {
		name          string
		input         *ListTournamentTeamsInput
		expectedErr   error
		expectedTeams int
		// expectedOwnedTeams int?
	}{
		{
			"empty",
			&ListTournamentTeamsInput{},
			ErrMissingArgument,
			0,
		}, {
			"unknown-tournament-id",
			&ListTournamentTeamsInput{TournamentID: "does not exist"},
			ErrInvalidArgument,
			0,
		}, {
			"solo-mode",
			&ListTournamentTeamsInput{TournamentID: tournaments["Solo Mode"]},
			nil,
			1,
		}, {
			"test-tournament",
			&ListTournamentTeamsInput{TournamentID: tournaments["Test Tournament"]},
			nil,
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.ListTournamentTeams(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			// fmt.Println(godev.PrettyJSON(ret))
			for _, team := range ret.Items {
				if team.TournamentID != test.input.TournamentID {
					t.Fatalf("Expected %q, got %q.", test.input.TournamentID, team.TournamentID)
				}
			}

			if len(ret.Items) != test.expectedTeams {
				t.Fatalf("Expected %d, got %d.", test.expectedTeams, len(ret.Items))
			}
		})
	}
}
