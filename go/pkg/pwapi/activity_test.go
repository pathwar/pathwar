package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
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
		// fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_UserRegister)
		assert.Equal(t, activity.Author.ID, session.User.ID)
		assert.Equal(t, activity.User.ID, session.User.ID)
		assert.Equal(t, activity.Team.ID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Global")
		assert.Equal(t, activity.Organization.ID, session.User.ActiveTeamMember.Team.Organization.ID)
		assert.Equal(t, activity.TeamMember.ID, session.User.ActiveTeamMember.ID)
	}

	// get session again
	{
		sess, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)
		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 1)
		assert.Equal(t, sess.User.ActiveTeamMember.Team.Cash, int64(0))
	}

	// FIXME: make a new login

	// FIXME: call UserSetPreferences

	// buy free challenge
	var freeChallenge, expensiveChallenge *pwdb.SeasonChallenge
	var subscription *SeasonChallengeBuy_Output
	{
		gs := testingGlobalSeason(t, svc)
		activeTeam := session.User.ActiveTeamMember.Team
		challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{gs.ID})
		require.NoError(t, err)

		for _, challenge := range challenges.Items {
			if challenge.Flavor.PurchasePrice == 0 && len(challenge.Flavor.Instances) > 0 {
				freeChallenge = challenge
			}
			if challenge.Flavor.PurchasePrice != 0 && len(challenge.Flavor.Instances) > 0 {
				expensiveChallenge = challenge
			}
		}
		require.NotNil(t, freeChallenge)
		require.NotNil(t, expensiveChallenge)
		input := SeasonChallengeBuy_Input{
			FlavorID: freeChallenge.Flavor.Slug,
			SeasonID: activeTeam.Season.Slug,
		}
		subscription, err = svc.SeasonChallengeBuy(ctx, &input)
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 2)
		activity := activities.Items[1]
		// fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_SeasonChallengeBuy)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Global")
		assert.Equal(t, activity.ChallengeSubscriptionID, subscription.ChallengeSubscription.ID)
		assert.Equal(t, activity.SeasonChallengeID, subscription.ChallengeSubscription.SeasonChallenge.ID)
	}

	// get session again
	{
		sess, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, sess.User.ActiveTeamMember.Team.Cash, int64(0))
	}

	// validate challenge
	{
		db := testingSvcDB(t, svc)
		// fetch full instance objects (base object is cleaned)
		err := db.First(&freeChallenge.Flavor.Instances[0], "ID = ?", freeChallenge.Flavor.Instances[0].ID).Error
		require.NoError(t, err)
		configData, err := freeChallenge.Flavor.Instances[0].ParseInstanceConfig()
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
		assert.Equal(t, activity.Season.Name, "Global")
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		// fmt.Println(godev.PrettyJSON(activity))
	}

	// get session again
	{
		sess, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, sess.User.ActiveTeamMember.Team.Cash, int64(10))
	}

	// validate coupon
	{
		input := CouponValidate_Input{
			Hash:   "test-coupon-1",
			TeamID: session.User.ActiveTeamMember.Team.ID,
		}
		ret, err := svc.CouponValidate(ctx, &input)
		require.NoError(t, err)
		// fmt.Println(godev.PrettyJSON(ret))

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 4)
		activity := activities.Items[3]
		// fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_CouponValidate)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Global")
		assert.Equal(t, activity.CouponID, ret.CouponValidation.CouponID)
	}

	// get session again
	{
		sess, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, sess.User.ActiveTeamMember.Team.Cash, int64(52))
	}

	// buy free challenge
	{
		activeTeam := session.User.ActiveTeamMember.Team

		subscription, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
			FlavorID: expensiveChallenge.Flavor.Slug,
			SeasonID: activeTeam.Season.Slug,
		})
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 5)
		activity := activities.Items[4]
		// fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_SeasonChallengeBuy)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.TeamID, session.User.ActiveTeamMember.Team.ID)
		assert.Equal(t, activity.Season.Name, "Global")
		assert.Equal(t, activity.ChallengeSubscriptionID, subscription.ChallengeSubscription.ID)
		assert.Equal(t, activity.SeasonChallengeID, subscription.ChallengeSubscription.SeasonChallenge.ID)
	}

	// get session again
	{
		sess, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, sess.User.ActiveTeamMember.Team.Cash, int64(47))
	}

	// delete account
	{
		_, err := svc.UserDeleteAccount(ctx, &UserDeleteAccount_Input{Reason: "testing activities"})
		require.NoError(t, err)

		activities = testingActivities(t, svc)
		require.Len(t, activities.Items, 6)
		activity := activities.Items[5]
		// fmt.Println(godev.PrettyJSON(activity))
		assert.Equal(t, activity.Kind, pwdb.Activity_UserDeleteAccount)
		assert.Equal(t, activity.AuthorID, session.User.ID)
		assert.Equal(t, activity.UserID, session.User.ID)
	}

	// get session again (re register)
	{
		sess, err := svc.UserGetSession(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, sess.User.ActiveTeamMember.Team.Cash, int64(0))
	}
}
