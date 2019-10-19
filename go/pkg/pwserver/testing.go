package pwserver

import (
	"context"
	"testing"

	"pathwar.land/go/pkg/pwengine"
)

func TestingServer(t *testing.T, ctx context.Context, opts Opts) (func() error, func()) {
	t.Helper()

	engine, engineCleanup := pwengine.TestingEngine(t, pwengine.Opts{Logger: opts.Logger})
	start, serverCleanup, err := Start(ctx, engine, opts)
	if err != nil {
		t.Fatalf("init server: %v", err)
	}

	cleanup := func() {
		serverCleanup()
		engineCleanup()
	}

	return start, cleanup
}
