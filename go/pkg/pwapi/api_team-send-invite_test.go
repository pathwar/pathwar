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

func TestService_TeamSendInvite(t *testing.T) {
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
	session2, err := svc.UserGetSession(ctx2, nil)
	require.NoError(t, err)

	seasonMap := map[string]*UserGetSession_Output_SeasonAndTeam{}
	for _, item := range session.Seasons {
		seasonMap[item.Season.Name] = item
	}

	var team1 pwdb.Team
	ret, err := svc.TeamCreate(ctx, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test1"].Season.ID),
		Name:     "Test1Team",
	})
	require.NoError(t, err)
	err = db.Model(pwdb.Team{}).Where(&pwdb.Team{ID: ret.Team.ID}).Update(&pwdb.Team{DeletionStatus: pwdb.DeletionStatus_Anonymized}).First(&team1).Error
	require.NoError(t, err)
	ret, err = svc.TeamCreate(ctx, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test1"].Season.ID),
		Name:     "Test2Team",
	})
	require.NoError(t, err)
	team2 := ret.Team
	ret, err = svc.TeamCreate(ctx, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test2"].Season.ID),
		Name:     "Test3Team",
	})
	require.NoError(t, err)
	team3 := ret.Team
	ret, err = svc.TeamCreate(ctx2, &TeamCreate_Input{
		SeasonID: fmt.Sprint(seasonMap["Test2"].Season.ID),
		Name:     "Test4Team",
	})
	require.NoError(t, err)
	team4 := ret.Team

	tests := []struct {
		name        string
		input       *TeamSendInvite_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &TeamSendInvite_Input{}, errcode.ErrMissingInput},
		{"only-team-id", &TeamSendInvite_Input{TeamID: fmt.Sprint(team1.ID)}, errcode.ErrMissingInput},
		{"invalid-team-id", &TeamSendInvite_Input{TeamID: "4242", UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrNoSuchSlug},
		{"deleted-team", &TeamSendInvite_Input{TeamID: fmt.Sprint(team1.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrTeamDoesNotExist},
		{"valid", &TeamSendInvite_Input{TeamID: fmt.Sprint(team2.ID), UserID: fmt.Sprint(session2.User.ID)}, nil},
		{"already-invited", &TeamSendInvite_Input{TeamID: fmt.Sprint(team2.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrAlreadyInvitedInTeam},
		{"has-team", &TeamSendInvite_Input{TeamID: fmt.Sprint(team3.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrAlreadyHasTeamForSeason},
		{"is-owner", &TeamSendInvite_Input{TeamID: fmt.Sprint(team4.ID), UserID: fmt.Sprint(session2.User.ID)}, errcode.ErrNotTeamOwner},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := svc.TeamSendInvite(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}
		})
	}
}
