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

func TestService_UserGetSession(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// register
	var session *UserGetSession_Output
	{
		var err error
		session, err = svc.UserGetSession(ctx, nil)
		require.NoError(t, err)

		// fmt.Println(godev.PrettyJSON(session))
		assert.Equal(t, session.User.Username, "mikael")
		assert.Len(t, session.Seasons, 3)
		assert.Equal(t, session.Claims.ActionToken.Sub, pwsso.TestingClaims(t).ActionToken.Sub)
		assert.True(t, session.IsNewUser)
		assert.Equal(t, session.User.ActiveTeamMember.Team.Season.Name, "Global")
		assert.Equal(t, session.User.ActiveTeamMember.Team.Organization.Name, "mikael")
		assert.True(t, session.User.ActiveTeamMember.Team.Organization.GlobalSeason)
		assert.Equal(t, session.User.ActiveTeamMember.Role, pwdb.TeamMember_Owner)
		for _, season := range session.Seasons {
			if season.Season.Name == "Global" {
				assert.Equal(t, season.Team.Organization.Name, "mikael")
			}
		}
	}

	// login
	{
		session2, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)

		// fmt.Println(godev.PrettyJSON(session2))
		assert.Equal(t, session2.User.Username, "mikael")
		assert.Len(t, session2.Seasons, 3)
		assert.Equal(t, session2.Claims.ActionToken.Sub, pwsso.TestingClaims(t).ActionToken.Sub)
		assert.False(t, session2.IsNewUser)
		assert.Equal(t, session2.User.ActiveTeamMember.Team.Season.Name, "Global")
		assert.Equal(t, session2.User.ActiveTeamMember.Team.Organization.Name, "mikael")
		assert.True(t, session2.User.ActiveTeamMember.Team.Organization.GlobalSeason)
		assert.Equal(t, session2.User.ActiveTeamMember.Role, pwdb.TeamMember_Owner)

		// standardize dynamic fields before comparison
		session.Notifications = session2.Notifications
		session.IsNewUser = false
		assert.Equal(t, session, session2)
	}
}
