package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func TestSvc_ChallengeBuy(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	gs := testingGlobalSeason(t, svc)

	// fetch user session
	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
	activeTeam := session.User.ActiveTeamMember.Team

	// fetch challenges
	challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{gs.ID})
	require.NoError(t, err)

	var expensiveChallenge, freeChallenge *pwdb.SeasonChallenge
	for _, challenge := range challenges.Items {
		if challenge.Flavor.PurchasePrice == 0 {
			freeChallenge = challenge
		} else {
			expensiveChallenge = challenge
		}
	}

	var tests = []struct {
		name        string
		input       *SeasonChallengeBuy_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &SeasonChallengeBuy_Input{}, errcode.ErrMissingInput},
		{"invalid flavor ID", &SeasonChallengeBuy_Input{FlavorID: "42", SeasonID: activeTeam.Season.Slug}, errcode.ErrInvalidFlavor},
		{"invalid season ID", &SeasonChallengeBuy_Input{FlavorID: freeChallenge.Flavor.Slug, SeasonID: "42"}, errcode.ErrInvalidSeason},
		{"not enough cash", &SeasonChallengeBuy_Input{FlavorID: expensiveChallenge.Flavor.Slug, SeasonID: activeTeam.Season.Slug}, errcode.ErrNotEnoughCash},
		{"valid 1", &SeasonChallengeBuy_Input{FlavorID: freeChallenge.Flavor.Slug, SeasonID: activeTeam.Season.Slug}, nil},
		{"valid 2 (duplicate)", &SeasonChallengeBuy_Input{FlavorID: freeChallenge.Flavor.Slug, SeasonID: activeTeam.Season.Slug}, errcode.ErrChallengeAlreadySubscribed},
		// FIXME: check for a team and a challenge in different seasons
		// FIXME: check for a team from another user
		// FIXME: check for a challenge in draft mode
	}

	for _, test := range tests {
		subscription, err := svc.SeasonChallengeBuy(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)
		if err != nil {
			continue
		}

		assert.Equalf(t, test.input.SeasonID, subscription.ChallengeSubscription.Team.Season.Slug, test.name)
		assert.Equalf(t, test.input.FlavorID, subscription.ChallengeSubscription.SeasonChallenge.Flavor.Slug, test.name)
		assert.Equalf(t, session.User.ID, subscription.ChallengeSubscription.BuyerID, test.name)

		// check if challenge subscription is now visible in season challenge list
		challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{gs.ID})
		if assert.NoError(t, err, test.name) {
			found := 0
			for _, challenge := range challenges.Items {
				if challenge.ID == subscription.ChallengeSubscription.SeasonChallengeID {
					found++
					if !assert.Lenf(t, challenge.Subscriptions, 1, test.name) {
						continue
					}
					assert.Equalf(t, subscription.ChallengeSubscription.ID, challenge.Subscriptions[0].ID, test.name)
					assert.Equalf(t, test.input.SeasonID, challenge.Subscriptions[0].Team.Season.Slug, test.name)
				}
			}
			assert.Equalf(t, 1, found, test.name)
		}
	}
}
