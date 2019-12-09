package pwagent

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"moul.io/godev"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwapi"
	"pathwar.land/go/pkg/pwcompose"
	"pathwar.land/go/pkg/pwdb"
)

func Daemon(ctx context.Context, cli *client.Client, logger *zap.Logger) error {
	// call API register in gRPC
	// ret, err := api.AgentRegister(ctx, &pwapi.AgentRegister_Input{Name: "dev", Hostname: "localhost", OS: "lorem ipsum", Arch: "x86_64", Version: "dev", Tags: []string{"dev"}})

	// list expected state from the output
	// apiInstances, err := api.AgentListInstances(ctx, &pwapi.AgentListInstances_Input{AgentID: ret.ID})
	apiInstances := &pwapi.AgentListInstances_Output{ // FIXME: tmp fake data; feel free to update it to match more cases
		Instances: []*pwdb.ChallengeInstance{
			{
				ID:             1,
				Status:         pwdb.ChallengeInstance_IsNew,
				InstanceConfig: []byte(`{"passphrases": [1, 2, 3, 4]}`),
				Flavor: &pwdb.ChallengeFlavor{
					ID:      2,
					Version: "latest",
					Challenge: &pwdb.Challenge{
						ID:   3,
						Name: "training-sqli",
					},
					SeasonChallenges: []*pwdb.SeasonChallenge{
						{
							ID: 4,
							Subscriptions: []*pwdb.ChallengeSubscription{
								{
									ID:     5,
									Status: pwdb.ChallengeSubscription_Active,
									Team: &pwdb.Team{
										ID: 6,
										Members: []*pwdb.TeamMember{
											{ID: 7},
											{ID: 8},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("api instances", godev.PrettyJSON(apiInstances))

	// fetch local info from docker daemon
	dockerInstances, err := pwcompose.GetPathwarInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetPathwarInfo.Wrap(err)
	}
	fmt.Println("info", godev.PrettyJSON(dockerInstances))

	// compute the difference

	// start missing instances

	// redump instances with state=NeedRedump

	// configure nginx
	// - compute the list of members in Active teams
	// - generate a hash based on their userIDs

	// later: for each updated instances -> call api to update status

	// last step -> make this a daemon: -> loop every X seconds

	// bonus: add optional parameters i.e. pathwar agent daemon --force-recreate

	return errcode.ErrNotImplemented
}
