package pwagent

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"moul.io/godev"
	"pathwar.land/v2/go/pkg/errcode"
	"pathwar.land/v2/go/pkg/pwapi"
	"pathwar.land/v2/go/pkg/pwcompose"
)

func updateAPIInstancesStatus(ctx context.Context, apiInstances *pwapi.AgentListInstances_Output, cli *client.Client, apiClient *pwapi.HTTPClient, opts Opts) error {
	//logger := opts.Logger

	containersInfo, err := pwcompose.GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	for _, flavor := range containersInfo.RunningFlavors {
		for _, apiInstance := range apiInstances.Instances {
			fmt.Println(flavor.ChallengeID(), godev.PrettyJSON(flavor), godev.PrettyJSON(apiInstance))
			//if flavor.ChallengeID() == apiInstance.Flavor.ChallengeID {
			//	apiInstance.Status = pwdb.ChallengeInstance_Available
			//}
		}
	}

	input := pwapi.AgentUpdateState_Input{
		Instances: apiInstances.Instances,
	}
	apiClient.AgentUpdateState(&input)

	return nil
}
