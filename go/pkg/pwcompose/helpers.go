package pwcompose

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func GetContainersInfo(ctx context.Context, cli *client.Client) (*ContainersInfo, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, errcode.ErrDockerAPIContainerList.Wrap(err)
	}

	containersInfo := ContainersInfo{
		RunningFlavors:    map[string]challengeFlavors{},
		RunningContainers: map[string]container{},
	}

	for _, dockerContainer := range containers {
		c := container(dockerContainer)

		// pathwar nginx proxy
		for _, name := range c.Names {
			if name[1:] == NginxContainerName {
				containersInfo.NginxContainer = c
			}
		}

		if _, found := c.Labels[challengeNameLabel]; !found { // not a pathwar container
			continue
		}

		flavor := c.ChallengeID()
		if _, found := containersInfo.RunningFlavors[flavor]; !found {
			challengeFlavor := challengeFlavors{
				Containers: map[string]container{},
			}
			challengeFlavor.Name = c.Labels[challengeNameLabel]
			challengeFlavor.Version = c.Labels[challengeVersionLabel]
			challengeFlavor.InstanceKey = c.Labels[InstanceKeyLabel]
			containersInfo.RunningFlavors[flavor] = challengeFlavor
		}
		containersInfo.RunningFlavors[flavor].Containers[c.ID] = c
		containersInfo.RunningContainers[c.ID] = c
	}

	// find proxy network
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, errcode.ErrDockerAPINetworkList.Wrap(err)
	}
	for _, networkResource := range networks {
		if networkResource.Name == ProxyNetworkName {
			containersInfo.NginxNetwork = networkResource
			break
		}
	}

	return &containersInfo, nil
}

func composeCliCommonArgs(path string) []string {
	return []string{"-f", path, "--no-ansi", "--log-level=ERROR"}
}
