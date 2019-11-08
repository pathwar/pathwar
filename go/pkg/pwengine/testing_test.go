package pwengine

import (
	"context"
	"testing"

	"github.com/jinzhu/gorm"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func testingSeasons(t *testing.T, e Engine) *pwdb.SeasonList {
	t.Helper()

	db := testingEngineDB(t, e)
	var list pwdb.SeasonList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list seasons: %v", err)
	}

	return &list
}

func testingSoloSeason(t *testing.T, e Engine) *pwdb.Season {
	t.Helper()

	seasons := testingSeasons(t, e)
	for _, season := range seasons.Items {
		if season.Name == "Solo Mode" {
			return season
		}
	}

	t.Fatalf("no such solo season")
	return nil
}

func testingTeams(t *testing.T, e Engine) *pwdb.TeamList {
	t.Helper()

	db := testingEngineDB(t, e)
	var list pwdb.TeamList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list season organizations: %v", err)
	}

	return &list
}

func testingChallenges(t *testing.T, e Engine) *pwdb.ChallengeList {
	t.Helper()

	db := testingEngineDB(t, e)
	var list pwdb.ChallengeList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list season organizations: %v", err)
	}

	return &list
}

func testingEngineDB(t *testing.T, e Engine) *gorm.DB {
	t.Helper()

	typed := e.(*engine)
	return typed.db
}

func testingSetContextToken(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	return context.WithValue(ctx, userTokenCtx, pwsso.TestingToken(t))
}
