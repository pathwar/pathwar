package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/v2/go/internal/testutil"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwdb"
)

func TestSvc_ChallengeSubscriptionClose(t *testing.T) {
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

	// buy two challenges
	subscription1, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[0].ID,
		TeamID:            activeTeam.ID,
	})
	require.NoError(t, err)
	subscription2, err := svc.SeasonChallengeBuy(ctx, &SeasonChallengeBuy_Input{
		SeasonChallengeID: challenges.Items[1].ID,
		TeamID:            activeTeam.ID,
	})
	require.NoError(t, err)

	// validate second challenge
	_, err = svc.ChallengeSubscriptionValidate(ctx, &ChallengeSubscriptionValidate_Input{
		ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID,
		Passphrase:              "secret",
	})
	require.NoError(t, err)

	var tests = []struct {
		name        string
		input       *ChallengeSubscriptionClose_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &ChallengeSubscriptionClose_Input{}, errcode.ErrMissingInput},
		{"subscription1", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription1.ChallengeSubscription.ID}, errcode.ErrMissingChallengeValidation},
		{"subscription2", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID}, nil},
		{"subscription2-again", &ChallengeSubscriptionClose_Input{ChallengeSubscriptionID: subscription2.ChallengeSubscription.ID}, errcode.ErrChallengeAlreadyClosed},
	}
	for _, test := range tests {
		ret, err := svc.ChallengeSubscriptionClose(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)
		if err != nil {
			continue
		}

		assert.NotNilf(t, ret.ChallengeSubscription.ClosedAt, test.name)
		assert.Equalf(t, session.User.ID, ret.ChallengeSubscription.CloserID, test.name)
		assert.Equalf(t, pwdb.ChallengeSubscription_Closed, ret.ChallengeSubscription.Status, test.name)
		assert.Equalf(t, activeTeam.ID, ret.ChallengeSubscription.Team.ID, test.name)
		assert.Equalf(t, test.input.ChallengeSubscriptionID, ret.ChallengeSubscription.ID, test.name)
		assert.NotEmptyf(t, ret.ChallengeSubscription.Validations, test.name)
	}
}
