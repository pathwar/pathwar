package pwapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func TestSvc_TeamAcceptInvite(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	db := testingSvcDB(t, svc)
	populate := []*pwdb.Season{
		{Name: "Test1", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		{Name: "Test2", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
	}
	for _, season := range populate {
		err := db.Create(season).Error
		require.NoError(t, err)
	}
	ctx := testingSetContextToken(context.Background(), t)
	session, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)

	ctx2 := testingSetContextToken2(context.Background(), t)
	_, err = svc.UserGetSession(ctx2, nil)
	require.NoError(t, err)

	seasonMap := map[string]*UserGetSession_Output_SeasonAndTeam{}
	for _, item := range session.Seasons {
		seasonMap[item.Season.Name] = item
	}

	ret, err := svc.TeamCreate(ctx2, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test1"].Season.ID),
		Name:     "Test1Team",
	})
	team1 := ret.Team
	ret2, err := svc.TeamSendInvite(ctx2, &TeamSendInvite_Input{
		TeamID: fmt.Sprint(team1.ID),
		UserID: fmt.Sprint(session.User.ID),
	})
	require.NoError(t, err)
	teamInvite1 := ret2.TeamInvite
	require.NoError(t, err)
	_, err = svc.TeamCreate(ctx, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test1"].Season.ID),
		Name:     "Test2Team",
	})
	require.NoError(t, err)
	ret, err = svc.TeamCreate(ctx2, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test2"].Season.ID),
		Name:     "Test3Team",
	})
	require.NoError(t, err)
	team2 := ret.Team
	ret2, err = svc.TeamSendInvite(ctx2, &TeamSendInvite_Input{
		TeamID: fmt.Sprint(team2.ID),
		UserID: fmt.Sprint(session.User.ID),
	})
	require.NoError(t, err)
	teamInvite2 := ret2.TeamInvite

	var tests = []struct {
		name        string
		input       *TeamAcceptInvite_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &TeamAcceptInvite_Input{}, errcode.ErrMissingInput},
		{"invalid-invite", &TeamAcceptInvite_Input{TeamInviteID: "yolo"}, errcode.ErrNoSuchSlug},
		{"has-team", &TeamAcceptInvite_Input{TeamInviteID: fmt.Sprint(teamInvite1.ID)}, errcode.ErrAlreadyHasTeamForSeason},
		{"valid", &TeamAcceptInvite_Input{TeamInviteID: fmt.Sprint(teamInvite2.ID)}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := svc.TeamAcceptInvite(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}
		})
	}
}
