package pwengine

import (
	"context"
	"testing"

	"moul.io/godev"
	"pathwar.land/go/pkg/pwsso"
)

func TestEngine_GetUserSession(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{})
	defer cleanup()
	ctx := testSetContextToken(t, context.Background())

	session, err := engine.GetUserSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//fmt.Println(godev.PrettyJSON(session))
	var tests = []struct {
		name     string
		actual   interface{}
		expected string
	}{
		{".User.Username", session.User.Username, `"moul"`},
		{"len(.Tournament)", len(session.Tournaments), "2"},
		{".Claims", session.Claims, godev.JSON(pwsso.TestingClaims(t))},
		{".IsNewUser", session.IsNewUser, `true`},
		{".User.ActiveTournamentMember.TournamentTeam.Tournament.Name", session.User.ActiveTournamentMember.TournamentTeam.Tournament.Name, `"Solo Mode"`},
		{".User.ActiveTournamentMember.TournamentTeam.Team.Name", session.User.ActiveTournamentMember.TournamentTeam.Team.Name, `"moul"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := godev.JSON(test.actual)
			if test.expected != actual {
				t.Fatalf("Expected: %q, got %q instead.", test.expected, actual)
			}
		})
	}
}
