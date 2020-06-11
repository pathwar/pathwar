package pwagent

import (
	"time"

	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/internal/randstring"
)

type Opts struct {
	DomainSuffix      string
	HostIP            string
	HostPort          string
	ModeratorPassword string
	AuthSalt          string
	ForceRecreate     bool
	NginxDockerImage  string
	Cleanup           bool
	RunOnce           bool
	LoopDelay         time.Duration
	DefaultAgent      bool
	Name              string
	NoRun             bool

	Logger *zap.Logger
}

func (opts *Opts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.AuthSalt == "" {
		opts.AuthSalt = randstring.RandString(10)
		opts.Logger.Warn("random salt generated", zap.String("salt", opts.AuthSalt))
	}
	if opts.ModeratorPassword == "" {
		opts.ModeratorPassword = randstring.RandString(10)
		opts.Logger.Warn("random moderator password generated", zap.String("password", opts.ModeratorPassword))
	}
}
