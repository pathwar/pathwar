package pwengine

import (
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/pwsso"
)

var _ Engine = (*engine)(nil)

type Engine interface {
	EngineServer

	Close() error
}

type engine struct {
	db        *gorm.DB
	opts      Opts
	sso       pwsso.Client
	startedAt time.Time
	logger    *zap.Logger
}

type Opts struct {
	Logger *zap.Logger
}

func New(db *gorm.DB, sso pwsso.Client, opts Opts) (Engine, error) {
	return &engine{
		logger:    opts.Logger,
		db:        db,
		opts:      opts,
		sso:       sso,
		startedAt: time.Now(),
	}, nil
}

func (e *engine) Close() error {
	// Note: everything passed in the New() should be closed by the parent of engine.
	// Here you need to close everything started by the engine itself.
	e.opts.Logger.Debug("closed engine")
	return nil
}
