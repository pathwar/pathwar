package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func TestActivity(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	activities := testingActivities(t, svc)
	assert.Len(t, activities.Items, 0)

	// register
	var session *UserGetSession_Output
	{
		var err error
		session, err = svc.UserGetSession(ctx, nil)
		assert.NoError(t, err)

		activities = testingActivities(t, svc)
		assert.Len(t, activities.Items, 1)
		activity := activities.Items[0]
		//fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_UserRegister)
		assert.Equal(t, activity.Author.ID, session.User.ID)
		assert.Equal(t, activity.User.ID, session.User.ID)
		assert.Equal(t, activity.Team.ID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Solo Mode")
		assert.Equal(t, activity.Organization.ID, session.User.ActiveTeamMember.Team.Organization.ID)
		assert.Equal(t, activity.TeamMember.ID, session.User.ActiveTeamMember.ID)
	}

	// login
	{
		session2, err := svc.UserGetSession(ctx, nil)
		assert.NoError(t, err)

		activities = testingActivities(t, svc)
		assert.Len(t, activities.Items, 2)
		activity := activities.Items[1]
		//fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_UserLogin)
		assert.Equal(t, activity.Author.ID, session2.User.ID)
		assert.Equal(t, activity.User.ID, session2.User.ID)
	}

	// delete account
	{
		_, err := svc.UserDeleteAccount(ctx, &UserDeleteAccount_Input{Reason: "testing activities"})
		assert.NoError(t, err)

		activities = testingActivities(t, svc)
		assert.Len(t, activities.Items, 3)
		activity := activities.Items[2]
		//fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_UserDeleteAccount)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.UserID, session.User.ID)
	}
}
