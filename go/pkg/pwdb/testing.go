package pwdb

import (
	"testing"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func TestingSqliteDB(t *testing.T, logger *zap.Logger) *gorm.DB {
	t.Helper()

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to initialize in-memory sqlite server: %v", err)
	}

	opts := Opts{
		Logger: logger,
		skipFK: true, // required for sqlite :(
	}
	db, err = Configure(db, opts)
	if err != nil {
		t.Fatalf("failed to initialize pwdb: %v", err)
	}

	return db
}

// FIXME: func TestingMySQLDB(t *testing.T, logger *zap.Logger) *gorm.DB { }
