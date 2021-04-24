package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/banner"
	"moul.io/motd"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwagent"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

func agentCommand() *ffcli.Command {
	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	agentFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	agentFlags.StringVar(&ssoOpts.ClientID, "sso-clientid", ssoOpts.ClientID, "SSO ClientID")
	agentFlags.StringVar(&ssoOpts.ClientSecret, "sso-clientsecret", ssoOpts.ClientSecret, "SSO ClientSecret")
	agentFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")
	agentFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")
	agentFlags.BoolVar(&agentOpts.Cleanup, "clean", agentOpts.Cleanup, "remove all pathwar instances before executing")
	agentFlags.BoolVar(&agentOpts.RunOnce, "once", agentOpts.RunOnce, "run once and don't start daemon loop")
	agentFlags.BoolVar(&agentOpts.NoRun, "no-run", agentOpts.NoRun, "stop after agent initialization (register and cleanup)")
	agentFlags.DurationVar(&agentOpts.LoopDelay, "delay", agentOpts.LoopDelay, "delay between each loop iteration")
	agentFlags.BoolVar(&agentOpts.DefaultAgent, "default-agent", agentOpts.DefaultAgent, "agent creates an instance for each available flavor on registration, else will only create an instance of debug-challenge")
	agentFlags.StringVar(&agentOpts.Name, "agent-name", agentOpts.Name, "Agent Name")
	agentFlags.StringVar(&agentOpts.DomainSuffix, "domain-suffix", agentOpts.DomainSuffix, "Domain suffix to append")
	agentFlags.StringVar(&agentOpts.NginxDockerImage, "docker-image", agentOpts.NginxDockerImage, "docker image used to generate nginx proxy container")
	agentFlags.StringVar(&agentOpts.HostIP, "host", agentOpts.HostIP, "Nginx HTTP listening addr")
	agentFlags.StringVar(&agentOpts.HostPort, "port", agentOpts.HostPort, "Nginx HTTP listening port")
	agentFlags.StringVar(&agentOpts.ModeratorPassword, "moderator-password", agentOpts.ModeratorPassword, "Challenge moderator password")
	agentFlags.StringVar(&agentOpts.AuthSalt, "salt", agentOpts.AuthSalt, "salt used to generate secure hashes (random if empty)")

	return &ffcli.Command{
		Name:       "agent",
		ShortUsage: "pathwar [global flags] agent [agent flags] <subcommand> [flags] [args...]",
		ShortHelp:  "manage an agent node (multiple challenges)",
		FlagSet:    agentFlags,
		Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
		Subcommands: []*ffcli.Command{
			{
				Name:      "pwinit.bin",
				ShortHelp: "dump pwinit binary to stdout",
				Exec: func(ctx context.Context, args []string) error {
					b, err := pwinit.Binary()
					if err != nil {
						return err
					}
					_, err = os.Stdout.Write(b)
					return err
				},
			},
		},
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("agent"))

			cleanup, err := initSentryFromEnv("starting agent")
			if err != nil {
				return err
			}
			defer cleanup()

			dockerCli, err := client.NewEnvClient()
			if err != nil {
				return errcode.ErrInitDockerClient.Wrap(err)
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			agentOpts.Logger = logger
			return pwagent.Run(ctx, dockerCli, apiClient, agentOpts)
		},
	}
}
