package pwdb

import (
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func TestingSqliteDB(t *testing.T, logger *zap.Logger) *gorm.DB {
	t.Helper()

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("init in-memory sqlite server: %v", err)
	}

	sfn, err := snowflake.NewNode(1)
	if err != nil {
		t.Fatalf("init snowflake generator: %v", err)
	}

	opts := Opts{
		Logger: logger,
		skipFK: true, // required for sqlite :(
	}

	db, err = Configure(db, sfn, opts)
	if err != nil {
		t.Fatalf("init pwdb: %v", err)
	}

	return db
}

// FIXME: func TestingMySQLDB(t *testing.T, logger *zap.Logger) *gorm.DB { }
