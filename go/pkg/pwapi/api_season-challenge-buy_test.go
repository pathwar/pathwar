package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/go/v2/internal/testutil"
	"pathwar.land/go/v2/pkg/errcode"
)

func TestSvc_ChallengeBuy(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	solo := testingSoloSeason(t, svc)

	// fetch user session
	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
	activeTeam := session.User.ActiveTeamMember.Team

	// fetch challenges
	challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
	require.NoError(t, err)

	var tests = []struct {
		name        string
		input       *SeasonChallengeBuy_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &SeasonChallengeBuy_Input{}, errcode.ErrMissingInput},
		{"invalid season challenge ID", &SeasonChallengeBuy_Input{SeasonChallengeID: 42, TeamID: activeTeam.ID}, errcode.ErrInvalidSeason},
		{"invalid team ID", &SeasonChallengeBuy_Input{SeasonChallengeID: challenges.Items[0].ID, TeamID: 42}, errcode.ErrInvalidTeam},
		{"valid 1", &SeasonChallengeBuy_Input{SeasonChallengeID: challenges.Items[0].ID, TeamID: activeTeam.ID}, nil},
		{"valid 2 (duplicate)", &SeasonChallengeBuy_Input{SeasonChallengeID: challenges.Items[0].ID, TeamID: activeTeam.ID}, errcode.ErrChallengeAlreadySubscribed},
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

		assert.Equalf(t, test.input.TeamID, subscription.ChallengeSubscription.TeamID, test.name)
		assert.Equalf(t, test.input.SeasonChallengeID, subscription.ChallengeSubscription.SeasonChallengeID, test.name)
		assert.Equalf(t, session.User.ID, subscription.ChallengeSubscription.BuyerID, test.name)

		// check if challenge subscription is now visible in season challenge list
		challenges, err := svc.SeasonChallengeList(ctx, &SeasonChallengeList_Input{solo.ID})
		if assert.NoError(t, err, test.name) {
			found := 0
			for _, challenge := range challenges.Items {
				if challenge.ID == subscription.ChallengeSubscription.SeasonChallengeID {
					found++
					if !assert.Lenf(t, challenge.Subscriptions, 1, test.name) {
						continue
					}
					assert.Equalf(t, subscription.ChallengeSubscription.ID, challenge.Subscriptions[0].ID, test.name)
					assert.Equalf(t, test.input.TeamID, challenge.Subscriptions[0].TeamID, test.name)
				}
			}
			assert.Equalf(t, 1, found, test.name)
		}
	}
}
