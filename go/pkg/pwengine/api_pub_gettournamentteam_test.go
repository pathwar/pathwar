package pwengine

import (
	"context"
	"errors"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_GetTournamentTeam(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// FIXME: check for permissions

	teams := map[string]int64{}
	for _, team := range testingTournamentTeams(t, engine).Items {
		teams[team.Team.Name] = team.ID
	}

	var tests = []struct {
		name                   string
		input                  *GetTournamentTeamInput
		expectedErr            error
		expectedTeamName       string
		expectedTournamentName string
	}{
		{
			"empty",
			&GetTournamentTeamInput{},
			ErrMissingArgument,
			"",
			"",
		}, {
			"unknown-tournament-id",
			&GetTournamentTeamInput{TournamentTeamID: -42}, // -42 should not exists
			ErrInvalidArgument,
			"",
			"",
		}, {
			"Staff",
			&GetTournamentTeamInput{TournamentTeamID: teams["Staff"]},
			nil,
			"Staff",
			"Solo Mode",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.GetTournamentTeam(ctx, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("Expected %#v, got %#v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}

			if ret.Item.ID != test.input.TournamentTeamID {
				t.Fatalf("Expected %q, got %q.", test.input.TournamentTeamID, ret.Item.ID)
			}
			if ret.Item.Team.Name != test.expectedTeamName {
				t.Fatalf("Expected %q, got %q.", test.expectedTeamName, ret.Item.Team.Name)
			}
			if ret.Item.Tournament.Name != test.expectedTournamentName {
				t.Fatalf("Expected %q, got %q.", test.expectedTournamentName, ret.Item.Tournament.Name)
			}
		})
	}
}
