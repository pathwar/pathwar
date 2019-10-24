package pwserver

import (
	"context"
	"testing"

	"pathwar.land/go/pkg/pwengine"
)

func TestingServer(t *testing.T, ctx context.Context, opts Opts) (*Server, func()) {
	t.Helper()

	engine, engineCleanup := pwengine.TestingEngine(t, pwengine.Opts{Logger: opts.Logger})

	if opts.HTTPBind == "" {
		opts.HTTPBind = "127.0.0.1:0"
	}
	if opts.GRPCBind == "" {
		opts.GRPCBind = "127.0.0.1:0"
	}

	server, err := New(ctx, engine, opts)
	if err != nil {
		t.Fatalf("init server: %v", err)
	}

	cleanup := func() {
		server.Close()
		engineCleanup()
	}

	go func() {
		if err := server.Run(); err != nil {
			t.Logf("server shutdown, err: %v", err)
		}
	}()

	return server, cleanup
}
