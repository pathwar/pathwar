package pwapi

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"pathwar.land/go/v2/pkg/pwdb"
	"pathwar.land/go/v2/pkg/pwsso"
)

func TestingService(t *testing.T, opts ServiceOpts) (Service, func()) {
	t.Helper()

	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	db := pwdb.TestingSqliteDB(t, opts.Logger)
	sso := pwsso.TestingSSO(t, opts.Logger)

	api, err := NewService(db, sso, opts)
	if err != nil {
		t.Fatalf("init api: %v", err)
	}

	cleanup := func() {
		api.Close()
		db.Close()
	}

	return api, cleanup
}

func TestingServer(t *testing.T, ctx context.Context, opts ServerOpts) (*Server, func()) {
	t.Helper()

	svc, svcCleanup := TestingService(t, ServiceOpts{Logger: opts.Logger})

	if opts.Bind == "" {
		opts.Bind = "127.0.0.1:0"
	}
	server, err := NewServer(ctx, svc, opts)
	assert.NoError(t, err)

	cleanup := func() {
		server.Close()
		svcCleanup()
	}

	go func() {
		if err := server.Run(); err != nil {
			opts.Logger.Warn("server shutdown", zap.Error(err))
		}
	}()

	return server, cleanup
}

func TestingClient(t *testing.T, address string) (ServiceClient, func()) {
	t.Helper()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	assert.NoError(t, err)
	c := NewServiceClient(conn)

	cleanup := func() {
		conn.Close()
	}
	return c, cleanup
}
