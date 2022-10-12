package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"moul.io/banner"
	"moul.io/motd"
	"os"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwagent"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"syscall"
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
		},
	}
}

func serverCommand() *ffcli.Command {
	devServerFlags := flag.NewFlagSet("dev", flag.ExitOnError)
	devServerFlags.StringVar(&serverOpts.Bind, "api-bind", ":8000", "api port (ex: :8000)")
	devServerFlags.BoolVar(&ssoOpts.AllowUnsafe, "sso-unsafe", true, "Allow unsafe SSO")
	devServerFlags.StringVar(&httpAPIAddr, "http-api-addr", "http://localhost:8000", "HTTP API address")
	devServerFlags.StringVar(&agentOpts.DomainSuffix, "domain-suffix", "localhost:8001", "Domain suffix to append")
	devServerFlags.BoolVar(&serverOpts.WithPprof, "with-pprof", true, "enable pprof endpoints")
	return &ffcli.Command{
		Name:      "server",
		ShortHelp: "launch api, agent & nginx",
		FlagSet:   devServerFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("server"))

			if err := globalPreRun(); err != nil {
				return err
			}

			cleanup, err := initSentryFromEnv("starting API")

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
				fmt.Println("Lancement API")
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
				fmt.Println("I'm here !")

				dockerCli, err := client.NewEnvClient()
				if err != nil {
					return errcode.ErrInitDockerClient.Wrap(err)
				}
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}
				server.Workers.Add(func() error {
					fmt.Println("Lancement Agent")
					err := pwagent.Run(ctx, dockerCli, apiClient, agentOpts)
					if err != cmux.ErrListenerClosed {
						return err
					}
					return nil
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

			return nil
		},
	}
}