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

func TestEngine_CouponValidate(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)
	db := testingSvcDB(t, svc)

	// fetch user session
	session, err := svc.UserGetSession(ctx, nil)
	assert.NoError(t, err)
	activeTeam := session.User.ActiveTeamMember.Team

	var tests = []struct {
		name         string
		input        *CouponValidate_Input
		expectedCash int64
		expectedErr  error
	}{
		{"nil", nil, 0, errcode.ErrMissingInput},
		{"empty", &CouponValidate_Input{}, 0, errcode.ErrMissingInput},
		{"invalid team ID", &CouponValidate_Input{Hash: "test-coupon-1", TeamID: 42}, 0, errcode.ErrUserDoesNotBelongToTeam},
		{"test-coupon-1", &CouponValidate_Input{Hash: "test-coupon-1", TeamID: activeTeam.ID}, 42, nil},
		{"test-coupon-1-again", &CouponValidate_Input{Hash: "test-coupon-1", TeamID: activeTeam.ID}, 42, errcode.ErrCouponAlreadyValidatedBySameTeam},
		{"test-coupon-2-invalid-season", &CouponValidate_Input{Hash: "test-coupon-2", TeamID: activeTeam.ID}, 42, errcode.ErrCouponNotFound},
		{"test-coupon-3-max-validation", &CouponValidate_Input{Hash: "test-coupon-3", TeamID: activeTeam.ID}, 42, errcode.ErrCouponExpired},
		{"test-coupon-4", &CouponValidate_Input{Hash: "test-coupon-4", TeamID: activeTeam.ID}, 84, nil},
	}

	for _, test := range tests {
		ret, err := svc.CouponValidate(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)

		// check cash, even if the previous function returned an error
		if test.input != nil && test.input.TeamID != 0 && test.input.TeamID != 42 {
			var team pwdb.Team
			err2 := db.First(&team, test.input.TeamID).Error
			require.NoErrorf(t, err2, test.name)
			assert.Equalf(t, test.expectedCash, team.Cash, test.name)
		}

		if err != nil {
			continue // skip other tests of previous function returned an error
		}

		validation := ret.CouponValidation
		assert.Equalf(t, test.input.Hash, validation.Coupon.Hash, test.name)
		assert.Equalf(t, test.input.TeamID, validation.Team.ID, test.name)
		assert.Equalf(t, validation.Team.SeasonID, validation.Coupon.SeasonID, test.name)
	}
}
