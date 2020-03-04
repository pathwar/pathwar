package pwagent

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"moul.io/godev"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwapi"
	"pathwar.land/v2/go/pkg/pwcompose"
	"pathwar.land/v2/go/pkg/pwdb"
)

func updateAPIState(ctx context.Context, apiInstances *pwapi.AgentListInstances_Output, cli *client.Client, apiClient *pwapi.HTTPClient, opts Opts) error {
	//logger := opts.Logger

	containersInfo, err := pwcompose.GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	for _, apiInstance := range apiInstances.Instances {
		for _, flavor := range containersInfo.RunningFlavors {
			isSame := false
			for _, container := range flavor.Containers {
				if container.Labels[pwcompose.InstanceKeyLabel] == fmt.Sprintf("%d", apiInstance.ID) {
					isSame = true
					break
				}
			}

			if isSame {
				// FIXME: check if state is "up"
				apiInstance.Status = pwdb.ChallengeInstance_Available
			}
		}

		// cleanup
		apiInstance.Flavor = nil
		apiInstance.Agent = nil
	}

	// FIXME: update state only if it changed
	input := pwapi.AgentUpdateState_Input{Instances: apiInstances.Instances}
	fmt.Println(godev.PrettyJSONPB(&input))
	// apiClient.AgentUpdateState(&input)

	return nil
}
