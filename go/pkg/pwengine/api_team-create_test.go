package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/pwdb"
)

func TestEngine_TeamCreate(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// populate more seasons
	db := testingEngineDB(t, engine)
	populate := []*pwdb.Season{
		&pwdb.Season{Name: "Test1", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		&pwdb.Season{Name: "Test2", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		&pwdb.Season{Name: "Test3", Status: pwdb.Season_Started, Visibility: pwdb.Season_Public},
		&pwdb.Season{Name: "Test4", Status: pwdb.Season_Stopped, Visibility: pwdb.Season_Public},
	}
	for _, season := range populate {
		err := db.Create(season).Error
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	session, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// create a non-solo organization
	nonSoloOrganization := pwdb.Organization{
		Name:       "non solo",
		SoloSeason: false,
		Members:    []*pwdb.OrganizationMember{{UserID: session.User.ID}},
	}
	err = db.Create(&nonSoloOrganization).Error
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	nonMemberOrganization := pwdb.Organization{
		Name:       "non member",
		SoloSeason: false,
	}
	err = db.Create(&nonMemberOrganization).Error
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	seasonMap := map[string]*UserGetSession_Output_SeasonAndTeam{}
	for _, item := range session.Seasons {
		seasonMap[item.Season.Name] = item
	}

	var tests = []struct {
		name        string
		input       *TeamCreate_Input
		expectedErr error
	}{
		{"nil", nil, ErrMissingArgument},
		{"empty", &TeamCreate_Input{}, ErrMissingArgument},
		{"only-season-id", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID}, ErrMissingArgument},
		{"invalid-season-id", &TeamCreate_Input{SeasonID: 4242, Name: "hello"}, ErrInvalidArgument},
		{"invalid-organization-id", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID, OrganizationID: 4242}, ErrInvalidArgument},
		{"blacklisted-name", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID, Name: " STAFF "}, ErrInvalidArgument},
		{"new-team-in-solo-mode-with-organization", &TeamCreate_Input{SeasonID: seasonMap["Solo Mode"].Season.ID, OrganizationID: session.User.ActiveTeamMember.Team.OrganizationID}, ErrInvalidArgument},
		{"new-team-in-solo-mode-with-name", &TeamCreate_Input{SeasonID: seasonMap["Solo Mode"].Season.ID, Name: "yolo"}, ErrInvalidArgument},
		{"too-many-arguments", &TeamCreate_Input{SeasonID: seasonMap["Solo Mode"].Season.ID, Name: "yolo", OrganizationID: session.User.ActiveTeamMember.Team.OrganizationID}, ErrInvalidArgument},
		{"conflict-org-name", &TeamCreate_Input{SeasonID: seasonMap["Test1"].Season.ID, Name: session.User.ActiveTeamMember.Team.Organization.Name}, ErrInvalidArgument},
		{"from-solo-organization", &TeamCreate_Input{SeasonID: seasonMap["Test2"].Season.ID, OrganizationID: session.User.ActiveTeamMember.Team.OrganizationID}, ErrInvalidArgument},
		{"non-member-organization", &TeamCreate_Input{SeasonID: seasonMap["Test2"].Season.ID, OrganizationID: nonMemberOrganization.ID}, ErrInvalidArgument},
		{"valid-with-organization", &TeamCreate_Input{SeasonID: seasonMap["Test2"].Season.ID, OrganizationID: nonSoloOrganization.ID}, nil},
		{"valid-with-name", &TeamCreate_Input{SeasonID: seasonMap["Test3"].Season.ID, Name: "yolo"}, nil},
		{"closed-season", &TeamCreate_Input{SeasonID: seasonMap["Test4"].Season.ID, Name: "yolo2"}, ErrInvalidArgument},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := engine.TeamCreate(ctx, test.input)
			if test.expectedErr != err {
				t.Errorf("Expected %v, got %v.", test.expectedErr, err)
			}
			if err != nil {
				return
			}
			//fmt.Println(godev.PrettyJSON(ret))

			if 1 != len(ret.Team.Members) {
				t.Errorf("Expected 1 team member, got %d.", len(ret.Team.Members))
			}
			if session.User.ID != ret.Team.Members[0].UserID {
				t.Errorf("Expected %d, got %d.", session.User.ID, ret.Team.Members[0].UserID)
			}
			if test.input.SeasonID != ret.Team.SeasonID {
				t.Errorf("Expected %d, got %d.", test.input.SeasonID, ret.Team.SeasonID)
			}
			if test.input.OrganizationID != 0 && test.input.OrganizationID != ret.Team.OrganizationID {
				t.Errorf("Expected %d, got %d.", test.input.OrganizationID, ret.Team.OrganizationID)
			}
			if ret.Team.Organization.SoloSeason {
				t.Errorf("Expected non-solo organization.")
			}
		})
	}
}
