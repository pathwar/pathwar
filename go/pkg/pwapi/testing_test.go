package pwapi

import (
	"context"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func testingSeasons(t *testing.T, svc Service) *pwdb.SeasonList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.SeasonList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list seasons")
	return &list
}

func testingAgents(t *testing.T, svc Service) *pwdb.AgentList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.AgentList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list agents: %v", err)
	}

	return &list
}

func testingSeasonChallenges(t *testing.T, svc Service) *pwdb.SeasonChallengeList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.SeasonChallengeList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list season challenges")
	return &list
}

func testingSoloSeason(t *testing.T, svc Service) *pwdb.Season {
	t.Helper()

	seasons := testingSeasons(t, svc)
	for _, season := range seasons.Items {
		if season.Name == "Solo Mode" {
			return season
		}
	}

	t.Fatalf("no such solo season")
	return nil
}

func testingTeams(t *testing.T, svc Service) *pwdb.TeamList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.TeamList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list teams")
	return &list
}

func testingChallenges(t *testing.T, svc Service) *pwdb.ChallengeList {
	t.Helper()

	db := testingSvcDB(t, svc)
	var list pwdb.ChallengeList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	assert.NoError(t, err, "list challenges")
	return &list
}

func testingSvcDB(t *testing.T, svc Service) *gorm.DB {
	t.Helper()

	typed := svc.(*service)
	return typed.db
}

func testingSetContextToken(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	return context.WithValue(ctx, userTokenCtx, pwsso.TestingToken(t))
}

func checkErr(t *testing.T, name string, err error) {
	t.Helper()

	if !assert.NoError(t, err) {
		t.Fatal(name)
	}
}

func testSameErrcodes(t *testing.T, name string, expected, got error) {
	t.Helper()

	prefix := ""
	if name != "" {
		prefix = name + ": "
	}
	assert.Equalf(
		t,
		errcode.ErrCode_name[errcode.Code(expected)],
		errcode.ErrCode_name[errcode.Code(got)],
		"%s%v", prefix, got,
	)
}
