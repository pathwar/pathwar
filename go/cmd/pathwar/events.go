package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/banner"
	"moul.io/motd"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwes"
)

func eventsCommand() *ffcli.Command {
	events := flag.NewFlagSet("events", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "events",
		ShortUsage: "pathwar [global flags] events [events flags] <cmd> [cmd flags]",
		ShortHelp:  "manage an event sourcing agent which process all events",
		FlagSet:    events,
		Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:       func(ctx context.Context, args []string) error { return flag.ErrHelp },
		Subcommands: []*ffcli.Command{
			eventsSourcing(),
			eventsRebuild(),
		},
	}
}

func eventsSourcing() *ffcli.Command {
	eventsSourcingFlags := flag.NewFlagSet("start", flag.ExitOnError)
	eventsSourcingFlags.IntVar(&esOpts.RefreshRate, "refresh-rate", esOpts.RefreshRate, "refresh rate in seconds")

	return &ffcli.Command{
		Name:      "start",
		ShortHelp: "start event sourcing agent",
		FlagSet:   eventsSourcingFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("event sourcing"))

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			var timestamp time.Time
			for {
				time.Sleep(time.Duration(esOpts.RefreshRate) * time.Second)
				err = pwes.EventHandler(ctx, apiClient, &timestamp, logger)
				if err != nil {
					return err
				}
			}
		},
	}
}

// TODO: Return error adapted
func eventsRebuild() *ffcli.Command {
	eventsRebuildFlags := flag.NewFlagSet("rebuild", flag.ExitOnError)
	eventsRebuildFlags.BoolVar(&esOpts.WithoutScore, "without-score", esOpts.WithoutScore, "rebuild without score")
	eventsRebuildFlags.StringVar(&esOpts.From, "from", esOpts.From, "rebuild from, format: YYYY-MM-DD HH:MM:SS")
	eventsRebuildFlags.StringVar(&esOpts.To, "to", esOpts.To, "rebuild to, format: YYYY-MM-DD HH:MM:SS")

	return &ffcli.Command{
		Name:      "rebuild",
		ShortHelp: "Rebuild current state from all events",
		FlagSet:   eventsRebuildFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("es rebuild"))

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

			err = pwes.Rebuild(ctx, apiClient, esOpts)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
