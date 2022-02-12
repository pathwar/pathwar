package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/pwversion"
)

func main() {
	err := runMain(os.Args[1:])
	switch {
	case err == nil:
		// noop
	case err == flag.ErrHelp || strings.Contains(err.Error(), flag.ErrHelp.Error()):
		os.Exit(2)
	default:
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}

func runMain(args []string) error {
	log.SetFlags(0)
	logger, _ := zap.NewProduction()
	defer func() {
		if logger != nil {
			_ = logger.Sync()
		}
	}()

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	var root *ffcli.Command
	{
		globalFlags := flag.NewFlagSet("pathwar", flag.ExitOnError)
		globalFlags.SetOutput(flagOutput) // used in main_test.go
		globalFlags.BoolVar(&globalDebug, "debug", false, "debug mode")
		globalFlags.StringVar(&zipkinEndpoint, "zipkin-endpoint", "", "optional opentracing server")
		globalFlags.StringVar(&bearerSecretKey, "bearer-secretkey", "", "bearer.sh secret key")
		globalFlags.StringVar(&globalSentryDSN, "sentry-dsn", "", "Sentry DSN")

		root = &ffcli.Command{
			ShortUsage: "pathwar [global flags] <subcommand> [flags] [args...]",
			FlagSet:    globalFlags,
			LongHelp:   "More info here: https://github.com/pathwar/pathwar/wiki/CLI",
			Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec:       func(ctx context.Context, args []string) error { return flag.ErrHelp },
			Subcommands: []*ffcli.Command{
				rawclientCommand(),
				cliCommand(),
				apiCommand(),
				composeCommand(),
				agentCommand(),
				miscCommand(),
				adminCommand(),
				{
					Name:       "version",
					ShortUsage: "pathwar [global flags] version",
					ShortHelp:  "show version",
					Exec: func(ctx context.Context, args []string) error {
						fmt.Printf(
							"version=%q\ncommit=%q\nbuilt-at=%q\nbuilt-by=%q\n",
							pwversion.Version, pwversion.Commit, pwversion.Date, pwversion.BuiltBy,
						)
						return nil
					},
				},
			},
		}
	}

	run := func() error {
		// create run.Group
		var process run.Group
		{
			// handle close signal
			execute, interrupt := run.SignalHandler(ctx, os.Interrupt)
			process.Add(execute, interrupt)

			// add root command to process
			process.Add(func() error {
				return root.ParseAndRun(ctx, args)
			}, func(error) {
				ctxCancel()
			})
		}

		// start the run.Group process
		{
			err := process.Run()
			if err == context.Canceled {
				return nil
			}
			return err
		}
	}

	return run()
}
