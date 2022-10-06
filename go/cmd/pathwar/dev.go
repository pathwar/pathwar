package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"moul.io/banner"
	"moul.io/motd"
	"os"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"syscall"
)

func devCommand() *ffcli.Command {
	cliFlags := flag.NewFlagSet("dev", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "dev",
		ShortUsage: "pathwar [global flags] dev [dev flags] <cmd> [cmd flags]",
		ShortHelp:  "carries out actions that help contribute to pathwar",
		FlagSet:    cliFlags,
		Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:       func(ctx context.Context, args []string) error { return flag.ErrHelp },
		Subcommands: []*ffcli.Command{
			serverCommand(),
		},
	}
}

func serverCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:      "server",
		ShortHelp: "start server (api + agent + nginx)",
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
				serverOpts.Tracer = tracer
				serverOpts.Logger = logger.Named("server")
				var err error

				if serverOpts.Bind == "gcloud" {
					serverOpts.Bind = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
					logger.Info("bind", zap.String("address", serverOpts.Bind))
				}

				server, err = pwapi.NewServer(ctx, svc, serverOpts)
				if err != nil {
					return errcode.ErrInitServer.Wrap(err)
				}
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
