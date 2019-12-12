package pwagent

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwapi"
	"pathwar.land/go/pkg/pwcompose"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwinit"
)

func Daemon(ctx context.Context, clean bool, runOnce bool, loopDelay time.Duration, cli *client.Client, logger *zap.Logger) error {
	// call API register in gRPC
	// ret, err := api.AgentRegister(ctx, &pwapi.AgentRegister_Input{Name: "dev", Hostname: "localhost", OS: "lorem ipsum", Arch: "x86_64", Version: "dev", Tags: []string{"dev"}})

	// list expected state from the output
	// apiInstances, err := api.AgentListInstances(ctx, &pwapi.AgentListInstances_Input{AgentID: ret.ID})
	apiInstances := &pwapi.AgentListInstances_Output{ // FIXME: tmp fake data; feel free to update it to match more cases
		Instances: []*pwdb.ChallengeInstance{
			{
				ID:             1,
				Status:         pwdb.ChallengeInstance_IsNew,
				InstanceConfig: []byte(`{"passphrases": ["a", "b", "c", "d"]}`),
				Flavor: &pwdb.ChallengeFlavor{
					ID:            2,
					Version:       "latest",
					ComposeBundle: "result of compose prepare",
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

	if clean {
		err := pwcompose.Down(ctx, []string{}, true, true, cli, logger)
		if err != nil {
			return errcode.ErrCleanPathwarInstances.Wrap(err)
		}
	}

	if runOnce {
		return run(ctx, apiInstances, cli, logger)
	}

	for {
		err := run(ctx, apiInstances, cli, logger)
		if err != nil {
			logger.Error("pwdaemon", zap.Error(err))
		}

		time.Sleep(loopDelay)
	}

	// later: for each updated instances -> call api to update status
}

func run(ctx context.Context, apiInstances *pwapi.AgentListInstances_Output, cli *client.Client, logger *zap.Logger) error {
	// fetch local info from docker daemon
	pathwarInfo, err := pwcompose.GetPathwarInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetPathwarInfo.Wrap(err)
	}

	agentOpts := AgentOpts{
		DomainSuffix:      "local",
		HostIP:            "0.0.0.0",
		HostPort:          "8000",
		ModeratorPassword: "",
		Salt:              "1337supmyman1337",
		AllowedUsers:      map[string][]int64{},
		ForceRecreate:     false,
		NginxDockerImage:  "docker.io/library/nginx:stable-alpine",
	}

	// compute instances that needs to upped / redumped
	for _, apiInstance := range apiInstances.GetInstances() {
		found := false
		needRedump := false
		for _, flavor := range pathwarInfo.RunningFlavors {
			if apiInstanceFlavor := apiInstance.GetFlavor(); apiInstanceFlavor != nil {
				if apiInstanceFlavorChallenge := apiInstanceFlavor.GetChallenge(); apiInstanceFlavorChallenge != nil {
					if flavor.InstanceKey == strconv.FormatInt(apiInstance.GetID(), 10) {
						found = true
						if apiInstance.GetStatus() == pwdb.ChallengeInstance_NeedRedump {
							needRedump = true
						}
					}
				}
			}
		}
		if !found || needRedump {
			// parse pwinit config
			var configData pwinit.InitConfig
			err = json.Unmarshal(apiInstance.GetInstanceConfig(), &configData)
			if err != nil {
				return errcode.ErrParseInitConfig.Wrap(err)
			}

			err = pwcompose.Up(ctx, apiInstance.GetFlavor().GetComposeBundle(), strconv.FormatInt(apiInstance.GetID(), 10), needRedump, &configData, cli, logger)
			if err != nil {
				return errcode.ErrUpPathwarInstance.Wrap(err)
			}
		}
	}

	// update pathwar infos
	pathwarInfo, err = pwcompose.GetPathwarInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetPathwarInfo.Wrap(err)
	}

	// update nginx configuration
	for _, apiInstance := range apiInstances.GetInstances() {
		if apiInstanceFlavor := apiInstance.GetFlavor(); apiInstanceFlavor != nil {
			if seasonChallenges := apiInstanceFlavor.GetSeasonChallenges(); seasonChallenges != nil {
				for _, seasonChallenge := range seasonChallenges {
					if subscriptions := seasonChallenge.GetActiveSubscriptions(); subscriptions != nil {
						for _, subscription := range subscriptions {
							if team := subscription.GetTeam(); team != nil {
								if members := team.GetMembers(); members != nil {
									for _, member := range members {
										for _, flavor := range pathwarInfo.RunningFlavors {
											if flavor.InstanceKey == strconv.FormatInt(apiInstance.GetID(), 10) {
												for _, instance := range flavor.Instances {
													for _, port := range instance.Ports {
														if port.PublicPort != 0 {
															// configure nginx
															// generate a hash per user for challenge dns prefix, based on their userIDs
															instanceName := instance.Names[0][1:]
															_, entryFound := agentOpts.AllowedUsers[instanceName]
															if !entryFound {
																agentOpts.AllowedUsers[instanceName] = []int64{member.GetID()}
															} else {
																allowedUsersSlice := agentOpts.AllowedUsers[instanceName]
																allowedUsersSlice = append(allowedUsersSlice, member.GetID())
																agentOpts.AllowedUsers[instanceName] = allowedUsersSlice
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	err = Nginx(ctx, agentOpts, cli, logger)
	if err != nil {
		return errcode.ErrUpdateNginx.Wrap(err)
	}

	return nil
}
