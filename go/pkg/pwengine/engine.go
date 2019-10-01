package pwengine

import (
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	pwsso "pathwar.land/go/pkg/pwsso"
)

var _ Client = (*client)(nil)

type Client interface {
	EngineServer

	Close() error
}

type client struct {
	db        *gorm.DB
	opts      Opts
	sso       pwsso.Client
	startedAt time.Time
	logger    *zap.Logger

	// implemented interfaces
	EngineServer
}

type Opts struct {
	Logger *zap.Logger
}

func New(db *gorm.DB, sso pwsso.Client, opts Opts) (Client, error) {
	return &client{
		logger:    opts.Logger,
		db:        db,
		opts:      opts,
		sso:       sso,
		startedAt: time.Now(),
	}, nil
}

func (c *client) Close() error {
	// Note: everything passed in the New() should be closed by the parent of engine.
	// Here you need to close everything started by the engine itself.
	c.opts.Logger.Debug("closed client")
	return nil
}
