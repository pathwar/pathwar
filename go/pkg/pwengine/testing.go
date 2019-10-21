package pwengine

import (
	"testing"

	"go.uber.org/zap"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

// TestEngine returns a configured Engine struct with in-memory contexts.
func TestingEngine(t *testing.T, opts Opts) (Engine, func()) {
	t.Helper()

	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	db := pwdb.TestingSqliteDB(t, opts.Logger)
	sso := pwsso.TestingSSO(t, opts.Logger)

	engine, err := New(db, sso, opts)
	if err != nil {
		t.Fatalf("init engine: %v", err)
	}

	cleanup := func() {
		engine.Close()
		db.Close()
	}

	return engine, cleanup
}
