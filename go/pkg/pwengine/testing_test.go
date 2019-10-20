package pwengine

import (
	"context"
	"testing"

	"github.com/jinzhu/gorm"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func testingTournaments(t *testing.T, e Engine) *pwdb.TournamentList {
	t.Helper()

	db := testingEngineDB(t, e)
	var list pwdb.TournamentList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list tournaments: %v", err)
	}

	return &list
}

func testingTournamentTeams(t *testing.T, e Engine) *pwdb.TournamentTeamList {
	t.Helper()

	db := testingEngineDB(t, e)
	var list pwdb.TournamentTeamList
	err := db.Set("gorm:auto_preload", true).Find(&list.Items).Error
	if err != nil {
		t.Fatalf("list tournament teams: %v", err)
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
