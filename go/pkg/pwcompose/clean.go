package pwcompose

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

// Purge cleans up everything related to Pathwar (containers, volumes, images, networks)
func Purge(ctx context.Context, cli *client.Client, logger *zap.Logger) error {
	return Clean(ctx, cli, CleanOpts{
		RemoveImages:  true,
		RemoveVolumes: true,
		RemoveNginx:   true,
		Logger:        logger,
	})
}

// DownAll cleans up everything related to Pathwar except images (containers, volumes, networks)
func DownAll(ctx context.Context, cli *client.Client, logger *zap.Logger) error {
	return Clean(ctx, cli, CleanOpts{
		RemoveVolumes: true,
		RemoveNginx:   true,
		Logger:        logger,
	})
}

type CleanOpts struct {
	ContainerIDs  []string
	RemoveImages  bool
	RemoveVolumes bool
	RemoveNginx   bool
	Logger        *zap.Logger
}

func NewCleanOpts() CleanOpts {
	return CleanOpts{
		RemoveVolumes: true,
	}
}

func (opts *CleanOpts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
}

// Clean can cleanup specific containers, all the images, all the volumes, and the pathwar's nginx front-end
func Clean(ctx context.Context, cli *client.Client, opts CleanOpts) error {
	opts.applyDefaults()
	opts.Logger.Debug("down", zap.Any("opts", opts))

	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	toRemove := map[string]container{}

	if opts.RemoveNginx && containersInfo.NginxContainer.ID != "" {
		toRemove[containersInfo.NginxContainer.ID] = containersInfo.NginxContainer
	}

	if len(opts.ContainerIDs) == 0 { // all containers
		for _, container := range containersInfo.RunningContainers {
			toRemove[container.ID] = container
		}
	} else { // only specific ones
		for _, id := range opts.ContainerIDs {
			for _, flavor := range containersInfo.RunningFlavors {
				if id == flavor.Name || id == flavor.ChallengeID() {
					for _, container := range flavor.Containers {
						toRemove[container.ID] = container
					}
				}
			}
			for _, container := range containersInfo.RunningContainers {
				if id == container.ID || id == container.ID[0:7] {
					toRemove[container.ID] = container
				}
			}
		}
	}

	for _, container := range toRemove {
		err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: opts.RemoveVolumes,
		})
		if err != nil {
			return errcode.ErrDockerAPIContainerRemove.Wrap(err)
		}
		opts.Logger.Debug("container removed", zap.String("ID", container.ID))
		if opts.RemoveImages {
			_, err := cli.ImageRemove(ctx, container.ImageID, types.ImageRemoveOptions{
				Force:         false,
				PruneChildren: true,
			})
			if err != nil {
				return errcode.ErrDockerAPIImageRemove.Wrap(err)
			}
			opts.Logger.Debug("image removed", zap.String("ID", container.ImageID))
		}
	}

	if opts.RemoveNginx && containersInfo.NginxNetwork.ID != "" {
		err = cli.NetworkRemove(ctx, containersInfo.NginxNetwork.ID)
		if err != nil {
			return errcode.ErrDockerAPINetworkRemove.Wrap(err)
		}
		opts.Logger.Debug("network removed", zap.String("ID", containersInfo.NginxNetwork.ID))
	}

	return nil
}
