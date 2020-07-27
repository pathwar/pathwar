package pwagent

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/client"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/internal/randstring"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
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

func Run(ctx context.Context, cli *client.Client, apiClient *pwapi.HTTPClient, opts Opts) error {
	started := time.Now()
	opts.applyDefaults()
	logger := opts.Logger

	err := agentRegister(ctx, apiClient, opts)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	if opts.Cleanup {
		before := time.Now()
		err := pwcompose.DownAll(ctx, cli, logger)
		if err != nil {
			return errcode.ErrCleanPathwarInstances.Wrap(err)
		}
		logger.Info("docker cleaned up", zap.Duration("duration", time.Since(before)))
	}

	if opts.NoRun {
		return nil
	}

	iteration := 0
	for {
		if !opts.RunOnce {
			logger.Debug("daemon iteration", zap.Int("number", iteration), zap.Duration("uptime", time.Since(started)))
		}

		err := runOnce(ctx, cli, apiClient, opts)
		if err != nil {
			logger.Error("daemon iteration", zap.Error(err))
		}

		if opts.RunOnce {
			break
		}

		opts.ForceRecreate = false // only do it once
		opts.Cleanup = false       // only do it once

		time.Sleep(opts.LoopDelay)
	}
	return nil
}

func runOnce(ctx context.Context, cli *client.Client, apiClient *pwapi.HTTPClient, opts Opts) error {
	instances, err := apiClient.AgentListInstances(ctx, &pwapi.AgentListInstances_Input{AgentName: opts.Name})
	opts.Logger.Debug("api response", zap.Any("instances", instances.GetInstances()))
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	if errs := applyDockerConfig(ctx, &instances, cli, opts); err != nil {
		for _, err := range multierr.Errors(errs) {
			opts.Logger.Error("apply docker config", zap.Error(err))
		}
	}

	if err := applyNginxConfig(ctx, &instances, cli, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := updateAPIState(ctx, &instances, cli, apiClient, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}

func NewOpts() Opts {
	return Opts{
		Cleanup:           false,
		RunOnce:           false,
		NoRun:             false,
		LoopDelay:         10 * time.Second,
		DefaultAgent:      true,
		Name:              getHostname(),
		DomainSuffix:      "local",
		NginxDockerImage:  "docker.io/library/nginx:stable-alpine",
		HostIP:            "0.0.0.0",
		HostPort:          "8001",
		ModeratorPassword: "",
		AuthSalt:          "",
	}
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

func getHostname() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		return "dev"
	}
	return hostname
}
