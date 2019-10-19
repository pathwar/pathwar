package pwengine

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
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
	snowflake *snowflake.Node
	startedAt time.Time
	logger    *zap.Logger
}

type Opts struct {
	Logger    *zap.Logger
	Snowflake *snowflake.Node
}

func New(db *gorm.DB, sso pwsso.Client, opts Opts) (Engine, error) {
	engine := &engine{
		logger:    opts.Logger,
		snowflake: opts.Snowflake,
		db:        db,
		opts:      opts,
		sso:       sso,
		startedAt: time.Now(),
	}

	if engine.logger == nil {
		engine.logger = zap.NewNop()
	}

	if engine.snowflake == nil {
		var err error
		engine.snowflake, err = snowflake.NewNode(1)
		if err != nil {
			return nil, fmt.Errorf("init snowflake: %w", err)
		}
	}

	return engine, nil
}

func (e *engine) Close() error {
	// Note: everything passed in the New() should be closed by the parent of engine.
	// Here you need to close everything started by the engine itself.
	e.opts.Logger.Debug("closed engine")
	return nil
}
