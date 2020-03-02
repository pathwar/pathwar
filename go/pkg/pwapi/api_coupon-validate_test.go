package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/v2/go/internal/testutil"
	"pathwar.land/v2/go/pkg/errcode"
)

func TestEngine_CouponValidate(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session
	session, err := svc.UserGetSession(ctx, nil)
	assert.NoError(t, err)
	activeTeam := session.User.ActiveTeamMember.Team

	var tests = []struct {
		name        string
		input       *CouponValidate_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &CouponValidate_Input{}, errcode.ErrMissingInput},
		// more invalid arguments
		// test-coupon-1, invalid season
		// {"invalid team ID", &CouponValidate_Input{CouponID: coupons.Items[0].ID, TeamID: 42}, ErrInvalidInput},
		// {"valid 1", &CouponValidate_Input{CouponID: coupons.Items[0].ID, TeamID: activeTeam.ID}, nil},
		// {"valid 2 (duplicate)", &CouponValidate_Input{CouponID: coupons.Items[0].ID, TeamID: activeTeam.ID}, errcode.TODO},
		// FIXME: check for a team and a coupon in different seasons
		// FIXME: check for a team from another user
		// FIXME: check for a coupon in draft mode
		{"test-coupon-1", &CouponValidate_Input{Hash: "test-coupon-1", TeamID: activeTeam.ID}, nil},
		{"test-coupon-1-again", &CouponValidate_Input{Hash: "test-coupon-1", TeamID: activeTeam.ID}, errcode.TODO},
		//{"test-coupon-2", &CouponValidate_Input{Hash: "test-coupon-2", TeamID: activeTeam.ID}, nil},
		{"test-coupon-3-invalid-season", &CouponValidate_Input{Hash: "test-coupon-3", TeamID: activeTeam.ID}, errcode.TODO},
		{"test-coupon-4", &CouponValidate_Input{Hash: "test-coupon-4", TeamID: activeTeam.ID}, nil},
	}

	for _, test := range tests {
		ret, err := svc.CouponValidate(ctx, test.input)
		testSameErrcodes(t, test.name, test.expectedErr, err)
		if err != nil {
			continue
		}

		validation := ret.CouponValidation
		assert.Equalf(t, test.input.Hash, validation.Coupon.Hash, test.name)
		assert.Equalf(t, test.input.TeamID, validation.Team.ID, test.name)
		assert.Equalf(t, validation.Team.SeasonID, validation.Coupon.SeasonID, test.name)
		//fmt.Println(godev.PrettyJSON(ret))
	}
}
