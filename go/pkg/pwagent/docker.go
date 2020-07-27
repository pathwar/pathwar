package pwagent

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func applyDockerConfig(ctx context.Context, apiInstances *pwapi.AgentListInstances_Output, dockerClient *client.Client, opts Opts) error {
	logger := opts.Logger
	logger.Debug("apply docker", zap.Any("opts", opts))

	// fetch local info from docker daemon
	containersInfo, err := pwcompose.GetContainersInfo(ctx, dockerClient)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	runningDockerInstances := map[string]bool{}
	for _, flavor := range containersInfo.RunningFlavors {
		runningDockerInstances[flavor.InstanceKey] = true
	}

	var (
		started        = 0
		ignored        = 0
		proxyNetworkID string
	)

	{ // configure proxy network
		networkResources, err := dockerClient.NetworkList(ctx, types.NetworkListOptions{})
		if err != nil {
			return errcode.ErrDockerAPINetworkList.Wrap(err)
		}
		for _, networkResource := range networkResources {
			if networkResource.Name == pwcompose.ProxyNetworkName {
				proxyNetworkID = networkResource.ID
				break
			}
		}
		if proxyNetworkID == "" {
			response, err := dockerClient.NetworkCreate(ctx, pwcompose.ProxyNetworkName, types.NetworkCreate{
				CheckDuplicate: true,
			})
			if err != nil {
				return errcode.ErrDockerAPINetworkCreate.Wrap(err)
			}
			proxyNetworkID = response.ID
			logger.Info("proxy network created", zap.String("name", pwcompose.ProxyNetworkName))
		}
	}

	var errs error
	for _, instance := range apiInstances.GetInstances() {
		instanceID := fmt.Sprintf("%d", instance.ID)
		l := logger.With(
			zap.String("id", instanceID),
			zap.String("flavor", instance.GetFlavor().NameAndVersion()),
		)
		isRunning := runningDockerInstances[instanceID]

		if instance.Status == pwdb.ChallengeInstance_Disabled {
			l.Debug("instance disabled")
			ignored++
			continue
		}

		if isRunning && instance.Status == pwdb.ChallengeInstance_Available {
			l.Debug("instance running")
			ignored++
			continue
		}

		// parse pwinit config
		configData, err := instance.ParseInstanceConfig()
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		bundle := instance.GetFlavor().GetComposeBundle()
		before := time.Now()
		started++

		upOpts := pwcompose.UpOpts{
			PreparedCompose: bundle,
			InstanceKey:     instanceID, // WARN -> normal?
			ForceRecreate:   true,
			ProxyNetworkID:  proxyNetworkID,
			PwinitConfig:    configData,
			Logger:          opts.Logger,
		}
		containers, err := pwcompose.Up(ctx, dockerClient, upOpts)
		if err != nil {
			errs = multierr.Append(errs, errcode.ErrUpPathwarInstance.Wrap(err))
			continue
		}

		l.Info(
			"started instance",
			zap.Duration("duration", time.Since(before)),
			zap.Int("containers", len(containers)),
		)
	}

	logger.Debug("docker stats", zap.Int("started", started), zap.Int("ignored", ignored), zap.Error(errs))

	return errs
}
