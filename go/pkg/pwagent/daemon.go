package pwagent

import (
	"context"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwcompose"
)

func Daemon(ctx context.Context, cli *client.Client, apiClient *http.Client, opts Opts) error {
	started := time.Now()

	err := opts.applyDefaults()
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	logger := opts.Logger

	// FIXME: call API register in gRPC
	// ret, err := api.AgentRegister(ctx, &pwapi.AgentRegister_Input{Name: "dev", Hostname: "localhost", OS: "lorem ipsum", Arch: "x86_64", Version: "dev", Tags: []string{"dev"}})

	if opts.Cleanup {
		before := time.Now()
		err := pwcompose.DownAll(ctx, cli, logger)
		if err != nil {
			return errcode.ErrCleanPathwarInstances.Wrap(err)
		}
		logger.Info("docker cleaned up", zap.Duration("duration", time.Since(before)))
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

		time.Sleep(opts.LoopDelay)
	}
	return nil
}

func runOnce(ctx context.Context, cli *client.Client, apiClient *http.Client, opts Opts) error {
	instances, err := fetchAPIInstances(ctx, apiClient, opts.HTTPAPIAddr, opts.Name, opts.Logger)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := applyDockerConfig(ctx, instances, cli, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := applyNginxConfig(ctx, instances, cli, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	// FIXME: implement this
	/*if err := updateAPIInstancesStatus(ctx, instances, cli, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}*/

	return nil
}
