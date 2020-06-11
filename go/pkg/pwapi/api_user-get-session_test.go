package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"moul.io/godev"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
)

func TestSvc_UserGetSession(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	//fmt.Println(godev.PrettyJSON(session))
	var tests = []struct {
		name     string
		actual   interface{}
		expected string
	}{
		{".User.Username", session.User.Username, `"moul"`},
		{"len(.Season)", len(session.Seasons), "2"},
		{".Claims", session.Claims, godev.JSON(pwsso.TestingClaims(t))},
		{".IsNewUser", session.IsNewUser, `true`},
		{".User.ActiveTeamMember.Team.Season.Name", session.User.ActiveTeamMember.Team.Season.Name, `"Solo Mode"`},
		{".User.ActiveTeamMember.Team.Organization.Name", session.User.ActiveTeamMember.Team.Organization.Name, `"moul"`},
		{".User.ActiveTeamMember.Team.Organization.SoloSeason", session.User.ActiveTeamMember.Team.Organization.SoloSeason, `true`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := godev.JSON(test.actual)
			assert.Equal(t, test.expected, actual)
		})
	}
}
