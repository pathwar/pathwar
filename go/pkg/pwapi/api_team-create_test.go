package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
)

func TestSvc_TeamCreate(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// populate more seasons
	db := testingSvcDB(t, svc)
	populate := []*pwdb.Season{
		&pwdb.Season{Name: "Test1", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		&pwdb.Season{Name: "Test2", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		&pwdb.Season{Name: "Test3", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		&pwdb.Season{Name: "Test4", Status: pwdb.Season_Stopped, Visibility: pwdb.Season_Public},
	}
	for _, season := range populate {
		err := db.Create(season).Error
		checkErr(t, "", err)
	}

	session, err := svc.UserGetSession(ctx, nil)
	checkErr(t, "", err)

	// create a non-solo organization
	nonSoloOrganization := pwdb.Organization{
		Name:       "non solo",
		SoloSeason: false,
		Members:    []*pwdb.OrganizationMember{{UserID: session.User.ID}},
	}
	err = db.Create(&nonSoloOrganization).Error
	checkErr(t, "", err)
	nonMemberOrganization := pwdb.Organization{
		Name:       "non member",
		SoloSeason: false,
	}
	err = db.Create(&nonMemberOrganization).Error
	checkErr(t, "", err)

	seasonMap := map[string]*UserGetSession_Output_SeasonAndTeam{}
	for _, item := range session.Seasons {
		seasonMap[item.Season.Name] = item
	}

	var tests = []struct {
		name        string
		input       *TeamCreate_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &TeamCreate_Input{}, errcode.ErrMissingInput},
		{"only-season-id", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID}, errcode.ErrMissingInput},
		{"invalid-season-id", &TeamCreate_Input{SeasonID: 4242, Name: "hello"}, errcode.ErrGetSeason},
		{"invalid-organization-id", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID, OrganizationID: 4242}, errcode.ErrGetOrganization},
		{"blacklisted-name", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID, Name: " STAFF "}, errcode.ErrReservedName},
		{"new-team-in-solo-mode-with-organization", &TeamCreate_Input{SeasonID: seasonMap["Solo Mode"].Season.ID, OrganizationID: session.User.ActiveTeamMember.Team.OrganizationID}, errcode.ErrAlreadyHasTeamForSeason},
		{"new-team-in-solo-mode-with-name", &TeamCreate_Input{SeasonID: seasonMap["Solo Mode"].Season.ID, Name: "yolo"}, errcode.ErrAlreadyHasTeamForSeason},
		{"too-many-arguments", &TeamCreate_Input{SeasonID: seasonMap["Solo Mode"].Season.ID, Name: "yolo", OrganizationID: session.User.ActiveTeamMember.Team.OrganizationID}, errcode.ErrInvalidInput},
		{"conflict-org-name", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID, Name: session.User.ActiveTeamMember.Team.Organization.Name}, errcode.ErrCheckOrganizationUniqueName},
		{"from-solo-organization", &TeamCreate_Input{SeasonID: seasonMap["Test2"].Season.ID, OrganizationID: session.User.ActiveTeamMember.Team.OrganizationID}, errcode.ErrCannotCreateTeamForSoloOrganization},
		{"non-member-organization", &TeamCreate_Input{SeasonID: seasonMap["Test2"].Season.ID, OrganizationID: nonMemberOrganization.ID}, errcode.ErrUserNotInOrganization},
		{"valid-with-organization", &TeamCreate_Input{SeasonID: seasonMap["Test2"].Season.ID, OrganizationID: nonSoloOrganization.ID}, nil},
		{"valid-with-name", &TeamCreate_Input{SeasonID: seasonMap["Test3"].Season.ID, Name: "yolo"}, nil},
		{"closed-season", &TeamCreate_Input{SeasonID: seasonMap["Test4"].Season.ID, Name: "yolo2"}, errcode.ErrSeasonDenied},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.TeamCreate(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Len(t, ret.Team.Members, 1)
			assert.Equal(t, session.User.ID, ret.Team.Members[0].UserID)
			assert.Equal(t, test.input.SeasonID, ret.Team.SeasonID)
			if test.input.OrganizationID != 0 {
				assert.Equal(t, test.input.OrganizationID, ret.Team.OrganizationID)
			}
			assert.False(t, ret.Team.Organization.SoloSeason)
		})
	}
}
