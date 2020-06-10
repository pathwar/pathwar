package pwagent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"moul.io/godev"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwapi"
	"pathwar.land/v2/go/pkg/pwcompose"
	"pathwar.land/v2/go/pkg/pwversion"
)

func Daemon(ctx context.Context, cli *client.Client, apiClient *pwapi.HTTPClient, opts Opts) error {
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

	if err := applyDockerConfig(ctx, &instances, cli, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := applyNginxConfig(ctx, &instances, cli, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	if err := updateAPIState(ctx, &instances, cli, apiClient, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}

func agentRegister(ctx context.Context, apiClient *pwapi.HTTPClient, opts Opts) error {
	metadata := map[string]interface{}{
		"delay":      opts.LoopDelay,
		"once":       opts.RunOnce,
		"num_cpus":   runtime.NumCPU(),
		"go_version": runtime.Version(),
		"uid":        os.Getuid(),
		"pid":        os.Getpid(),
	}
	metadataStr, err := json.Marshal(metadata)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	nginxPort, _ := strconv.Atoi(opts.HostPort)
	ret, err := apiClient.AgentRegister(ctx, &pwapi.AgentRegister_Input{
		Name:         opts.Name,
		Hostname:     hostname,
		NginxPort:    int32(nginxPort),
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		Version:      pwversion.Version,
		Tags:         []string{},
		DomainSuffix: opts.DomainSuffix,
		AuthSalt:     opts.AuthSalt,
		Metadata:     string(metadataStr),
		DefaultAgent: opts.DefaultAgent,
	})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		fmt.Fprintln(os.Stderr, godev.PrettyJSON(ret))
	}
	return nil
}
