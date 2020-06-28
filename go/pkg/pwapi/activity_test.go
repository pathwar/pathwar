package pwapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

func TestActivity(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	activities := testingActivities(t, svc)
	require.Len(t, activities.Items, 0)

	// register
	var session *UserGetSession_Output
	{
		var err error
		session, err = svc.UserGetSession(ctx, nil)
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 1)
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

	// get session again
	{
		_, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 1)
	}

	// FIXME: make a new login

	// FIXME: call UserSetPreferences

	// buy challenge
	var theChallenge *pwdb.SeasonChallenge
	var subscription *SeasonChallengeBuy_Output
	{
		solo := testingSoloSeason(t, svc)
		activeTeam := session.User.ActiveTeamMember.Team
		challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
		require.NoError(t, err)

		for _, challenge := range challenges.Items {
			if len(challenge.Flavor.Instances) > 0 {
				theChallenge = challenge
				break
			}
		}
		require.NotNil(t, theChallenge)
		subscription, err = svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{SeasonChallengeID: theChallenge.ID, TeamID: activeTeam.ID})
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 2)
		activity := activities.Items[1]
		//fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_SeasonChallengeBuy)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Solo Mode")
		assert.Equal(t, activity.ChallengeSubscriptionID, subscription.ChallengeSubscription.ID)
		assert.Equal(t, activity.SeasonChallengeID, subscription.ChallengeSubscription.SeasonChallenge.ID)
	}

	// validate challenge
	{
		var configData pwinit.InitConfig
		err := json.Unmarshal(theChallenge.Flavor.Instances[0].GetInstanceConfig(), &configData)
		require.NoError(t, err)
		input := ChallengeSubscriptionValidate_Input{
			ChallengeSubscriptionID: subscription.ChallengeSubscription.ID,
			Passphrases:             configData.Passphrases,
		}
		_, err = svc.ChallengeSubscriptionValidate(ctx, &input)
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 3)
		activity := activities.Items[2]
		assert.Equal(t, activity.Kind, pwdb.Activity_ChallengeSubscriptionValidate)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.ChallengeSubscriptionID, subscription.ChallengeSubscription.ID)
		assert.Equal(t, activity.SeasonChallengeID, subscription.ChallengeSubscription.SeasonChallenge.ID)
		assert.Equal(t, activity.ChallengeFlavorID, subscription.ChallengeSubscription.SeasonChallenge.Flavor.ID)
		assert.Equal(t, activity.Season.Name, "Solo Mode")
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		//fmt.Println(godev.PrettyJSON(activity))
	}

	// validate coupon
	{
		input := CouponValidate_Input{
			Hash:   "test-coupon-1",
			TeamID: session.User.ActiveTeamMember.Team.ID,
		}
		ret, err := svc.CouponValidate(ctx, &input)
		require.NoError(t, err)
		//fmt.Println(godev.PrettyJSON(ret))

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 4)
		activity := activities.Items[3]
		//fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_CouponValidate)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Solo Mode")
		assert.Equal(t, activity.CouponID, ret.CouponValidation.CouponID)
	}

	// delete account
	{
		_, err := svc.UserDeleteAccount(ctx, &UserDeleteAccount_Input{Reason: "testing activities"})
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 5)
		activity := activities.Items[4]
		//fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_UserDeleteAccount)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.UserID, session.User.ID)
	}
}
