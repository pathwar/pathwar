package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
)

func TestSvc_UserGetSession(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	// fmt.Println(godev.PrettyJSON(session))
	assert.Equal(t, session.User.Username, "moul")
	assert.Len(t, session.Seasons, 2)
	assert.Equal(t, session.Claims, pwsso.TestingClaims(t))
	assert.True(t, session.IsNewUser)
	assert.Equal(t, session.User.ActiveTeamMember.Team.Season.Name, "Solo Mode")
	assert.Equal(t, session.User.ActiveTeamMember.Team.Organization.Name, "moul")
	assert.True(t, session.User.ActiveTeamMember.Team.Organization.SoloSeason)
	assert.Equal(t, session.User.ActiveTeamMember.Role, pwdb.TeamMember_Owner)
	for _, season := range session.Seasons {
		if season.Season.Name == "Solo Mode" {
			assert.Equal(t, season.Team.Organization.Name, "moul")
		}
	}
}
