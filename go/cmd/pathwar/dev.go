package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/pwes"

	"github.com/docker/docker/client"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"moul.io/banner"
	"moul.io/motd"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwagent"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
)

func devCommand() *ffcli.Command {
	devFlags := flag.NewFlagSet("dev", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "dev",
		ShortUsage: "pathwar [global flags] dev [dev flags] <cmd> [cmd flags]",
		ShortHelp:  "carries out actions that help contribute to pathwar",
		FlagSet:    devFlags,
		Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:       func(ctx context.Context, args []string) error { return flag.ErrHelp },
		Subcommands: []*ffcli.Command{
			serverCommand(),
			challengeRunCommand(),
			challengeDeployCommand(),
			rollbackScore(),
		},
	}
}

func serverCommand() *ffcli.Command {
	devServerFlags := flag.NewFlagSet("dev", flag.ExitOnError)
	devServerFlags.StringVar(&serverOpts.Bind, "api-bind", ":8000", "api port (ex: :8000)")
	devServerFlags.BoolVar(&ssoOpts.AllowUnsafe, "sso-unsafe", true, "Allow unsafe SSO")
	devServerFlags.StringVar(&httpAPIAddr, "http-api-addr", "http://localhost:8000", "HTTP API address")
	devServerFlags.StringVar(&agentOpts.HostPort, "host-port", "8001", "Port nginx")
	devServerFlags.StringVar(&agentOpts.DomainSuffix, "domaine-suffix", "localhost:8001", "Domain suffix to append")
	devServerFlags.BoolVar(&serverOpts.WithPprof, "with-pprof", true, "enable pprof endpoints")

	return &ffcli.Command{
		Name:      "server",
		ShortHelp: "launch api, agent & nginx",
		FlagSet:   devServerFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())

			if err := globalPreRun(); err != nil {
				return err
			}

			cleanup, err := initSentryFromEnv("starting API")
			if err != nil {
				return err
			}

			svc, _, closer, err := svcFromFlags(logger)
			if err != nil {
				return errcode.ErrStartService.Wrap(err)
			}
			defer closer()

			if err != nil {
				return err
			}
			defer cleanup()

			var (
				g      run.Group
				server *pwapi.Server
			)

			g.Add(run.SignalHandler(ctx, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill))
			{
				serverOpts.Tracer = tracer
				serverOpts.Logger = logger.Named("server")
				fmt.Println(serverOpts)
				var err error

				if serverOpts.Bind == "gcloud" {
					serverOpts.Bind = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
					logger.Info("bind", zap.String("address", serverOpts.Bind))
				}

				server, err = pwapi.NewServer(ctx, svc, serverOpts)
				if err != nil {
					return errcode.ErrInitServer.Wrap(err)
				}

				dockerCli, err := client.NewEnvClient()
				if err != nil {
					return errcode.ErrInitDockerClient.Wrap(err)
				}
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				server.Workers.Add(func() error {
					err := pwagent.Run(ctx, dockerCli, apiClient, agentOpts)
					if err != cmux.ErrListenerClosed {
						return err
					}
					return nil
				}, func(error) {
					_, cancel := context.WithTimeout(ctx, 5)
					defer cancel()
				})

				server.Workers.Add(func() error {
					if err != nil {
						return err
					}
					for {
						time.Sleep(10 * time.Second)
						err = pwes.RollbackScore(ctx, apiClient)
						if err != nil {
							return err
						}
					}
				}, func(error) {
					_, cancel := context.WithTimeout(ctx, 5)
					defer cancel()
				})

				g.Add(
					server.Run,
					func(error) { server.Close() },
				)
			}

			logger.Info("server started",
				zap.String("bind", server.ListenerAddr()),
			)

			if err := g.Run(); err != nil {
				return errcode.ErrGroupTerminated.Wrap(err)
			}
			return nil
		},
	}
}

func challengeRunCommand() *ffcli.Command {
	devChallengeFlags := flag.NewFlagSet("dev", flag.ExitOnError)

	var (
		composePrepareOpts = pwcompose.NewPrepareOpts()
	)

	return &ffcli.Command{
		Name:      "challenge-run",
		ShortHelp: "register a challenge",
		FlagSet:   devChallengeFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("run challenge"))

			if err := globalPreRun(); err != nil {
				return err
			}

			composePrepareOpts.ChallengeDir = "."
			composePrepareOpts.Logger = logger
			composePrepareOpts.NoPush = true

			preparedComposeData, err := pwcompose.Prepare(composePrepareOpts)
			if err != nil {
				return err
			}

			var config pwcompose.PathwarConfig
			if err = yaml.Unmarshal([]byte(preparedComposeData), &config); err != nil {
				return errcode.TODO.Wrap(err)
			}
			config.Pathwar.Flavor.ComposeBundle = preparedComposeData

			slug := config.Pathwar.Challenge.Slug
			if slug == "" {
				return errors.New("a challenge slug is required in docker-compose.yml")
			}

			input := pwapi.AdminChallengeAdd_Input{
				Challenge: &config.Pathwar.Challenge,
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			_, err = apiClient.AdminAddChallenge(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			flavor := pwapi.AdminChallengeFlavorAdd_Input{
				ChallengeID:     slug,
				ChallengeFlavor: &config.Pathwar.Flavor,
			}

			_, err = apiClient.AdminAddChallengeFlavor(ctx, &flavor)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			season := pwapi.AdminSeasonChallengeAdd_Input{
				FlavorID: slug,
				SeasonID: "global",
			}

			_, err = apiClient.AdminAddSeasonChallenge(ctx, &season)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			_, err = apiClient.AdminRedumpChallenge(ctx, &pwapi.AdminChallengeRedump_Input{
				ChallengeID: slug,
			})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			return nil
		},
	}
}

func challengeDeployCommand() *ffcli.Command {
	devChallengeFlags := flag.NewFlagSet("dev", flag.ExitOnError)

	return &ffcli.Command{
		Name:      "challenge-deploy",
		ShortHelp: "deploy a challenge",
		FlagSet:   devChallengeFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("deploy challenge"))

			//TODO: Deploying challenge from current folder

			return nil
		},
	}
}

// TODO: Return error adapted
func rollbackScore() *ffcli.Command {
	var (
		devComputeFlags = flag.NewFlagSet("compute-score", flag.ExitOnError)
	)

	return &ffcli.Command{
		Name:      "rollback-score",
		ShortHelp: "Compute the score thanks to retrieving all events",
		FlagSet:   devComputeFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("compute score"))

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if err != nil {
				return err
			}

			err = pwes.RollbackScore(ctx, apiClient)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
