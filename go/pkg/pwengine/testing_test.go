package pwengine

import (
	"testing"

	"github.com/jinzhu/gorm"
	"pathwar.land/go/pkg/pwdb"
)

func testingTournaments(t *testing.T, e Engine) *pwdb.TournamentList {
	t.Helper()

	db := testingEngineDB(t, e)
	var list pwdb.TournamentList
	err := db.Find(&list.Items).Error
	if err != nil {
		t.Fatalf("failed to list tournaments: %v", err)
	}

	return &list
}

func testingEngineDB(t *testing.T, e Engine) *gorm.DB {
	t.Helper()

	typed := e.(*engine)
	return typed.db
}
