package pwagent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwapi"
	"pathwar.land/go/pkg/pwcompose"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwinit"
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

	for _, instance := range apiInstances.GetInstances() {
		instanceID := fmt.Sprintf("%d", instance.ID)
		l := logger.With(
			zap.String("id", instanceID),
			zap.String("flavor", instance.GetFlavor().NameAndVersion()),
		)
		if instance.Status == pwdb.ChallengeInstance_Disabled {
			l.Debug("instance disabled")
			ignored++
			continue
		}

		isRunning := runningDockerInstances[instanceID]
		if isRunning && instance.Status != pwdb.ChallengeInstance_NeedRedump {
			l.Debug("instance running")
			ignored++
			continue
		}

		// parse pwinit config
		var configData pwinit.InitConfig
		err = json.Unmarshal(instance.GetInstanceConfig(), &configData)
		if err != nil {
			return errcode.ErrParseInitConfig.Wrap(err)
		}

		bundle := instance.GetFlavor().GetComposeBundle()
		before := time.Now()
		started++
		containers, err := pwcompose.Up(ctx, bundle, instanceID, true, proxyNetworkID, &configData, dockerClient, logger)
		if err != nil {
			return errcode.ErrUpPathwarInstance.Wrap(err)
		}

		l.Info(
			"started instance",
			zap.Duration("duration", time.Since(before)),
			zap.Int("containers", len(containers)),
		)
	}

	logger.Debug("docker stats", zap.Int("started", started), zap.Int("ignored", ignored))

	return nil
}
