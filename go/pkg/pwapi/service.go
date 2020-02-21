package pwapi

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"pathwar.land/go/v2/pkg/errcode"
	"pathwar.land/go/v2/pkg/pwsso"
)

type Service interface {
	ServiceServer // generated with protobuf

	Close() error
}

func NewService(db *gorm.DB, sso pwsso.Client, opts ServiceOpts) (Service, error) {
	svc := &service{
		logger:    opts.Logger,
		snowflake: opts.Snowflake,
		db:        db,
		opts:      opts,
		sso:       sso,
		startedAt: time.Now(),
	}

	if svc.logger == nil {
		svc.logger = zap.NewNop()
	}

	if svc.snowflake == nil {
		var err error
		svc.snowflake, err = snowflake.NewNode(1)
		if err != nil {
			return nil, errcode.ErrInitSnowflake.Wrap(err)
		}
	}

	return svc, nil
}

type service struct {
	db        *gorm.DB
	opts      ServiceOpts
	sso       pwsso.Client
	snowflake *snowflake.Node
	startedAt time.Time
	logger    *zap.Logger
}

type ServiceOpts struct {
	Logger    *zap.Logger
	Snowflake *snowflake.Node
}

func (svc *service) Close() error {
	// Note: everything passed in the New() should be closed by the parent of service.
	// Here you need to close everything started by the service itself.
	svc.opts.Logger.Debug("closed service")
	return nil
}

var _ Service = (*service)(nil)
