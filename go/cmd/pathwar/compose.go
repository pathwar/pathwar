package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/client"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
)

func composeCommand() *ffcli.Command {
	var (
		composeDownFlags    = flag.NewFlagSet("compose down", flag.ExitOnError)
		composeFlags        = flag.NewFlagSet("compose", flag.ExitOnError)
		composePSFlags      = flag.NewFlagSet("compose ps", flag.ExitOnError)
		composePrepareFlags = flag.NewFlagSet("compose prepare", flag.ExitOnError)
		composeUpFlags      = flag.NewFlagSet("compose up", flag.ExitOnError)

		composeCleanOpts   = pwcompose.NewCleanOpts()
		composePrepareOpts = pwcompose.NewPrepareOpts()
		composeUpOpts      = pwcompose.NewUpOpts()
	)
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveVolumes, "rm-volumes", composeCleanOpts.RemoveVolumes, "keep volumes")
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveImages, "rm-images", composeCleanOpts.RemoveImages, "remove images as well")
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveNginx, "rm-nginx", composeCleanOpts.RemoveNginx, "down nginx container and proxy network as well")
	composePSFlags.IntVar(&composePSDepth, "depth", 0, "depth to display")
	composePrepareFlags.BoolVar(&composePrepareOpts.NoPush, "no-push", composePrepareOpts.NoPush, "don't push images")
	composePrepareFlags.StringVar(&composePrepareOpts.Prefix, "prefix", composePrepareOpts.Prefix, "docker image prefix")
	composePrepareFlags.StringVar(&composePrepareOpts.Version, "version", composePrepareOpts.Version, "challenge version")
	composePrepareFlags.BoolVar(&composePrepareOpts.JSON, "json", composePrepareOpts.JSON, "JSON format")
	composeUpFlags.StringVar(&composeUpOpts.InstanceKey, "instance-key", composeUpOpts.InstanceKey, "instance key used to generate instance ID")
	composeUpFlags.BoolVar(&composeUpOpts.ForceRecreate, "force-recreate", composeUpOpts.ForceRecreate, "down previously created instances of challenge")

	return &ffcli.Command{
		Name:  "compose",
		Usage: "pathwar [global flags] compose [compose flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{{
			Name:    "up",
			Usage:   "pathwar [global flags] compose [compose flags] up [flags] PATH",
			FlagSet: composeUpFlags,
			Options: []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(args []string) error {
				if len(args) != 1 {
					return flag.ErrHelp
				}

				err := globalPreRun()
				if err != nil {
					return err
				}

				path := args[0]
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				preparedCompose, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}

				ctx := context.Background()
				cli, err := client.NewEnvClient()
				if err != nil {
					return errcode.ErrInitDockerClient.Wrap(err)
				}

				composeUpOpts.Logger = logger
				composeUpOpts.PreparedCompose = string(preparedCompose)
				services, err := pwcompose.Up(ctx, cli, composeUpOpts)
				if err != nil {
					return err
				}

				for _, service := range services {
					fmt.Println(service.ContainerName)
				}

				return nil
			},
		}, {
			Name:    "prepare",
			Usage:   "pathwar [global flags] compose [compose flags] prepare [flags] PATH",
			FlagSet: composePrepareFlags,
			Options: []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(args []string) error {
				if len(args) != 1 {
					return flag.ErrHelp
				}
				path := args[0]
				err := globalPreRun()
				if err != nil {
					return err
				}

				composePrepareOpts.ChallengeDir = path
				composePrepareOpts.Logger = logger
				preparedComposeData, err := pwcompose.Prepare(composePrepareOpts)
				fmt.Println(preparedComposeData)
				return err
			},
		}, {
			Name:    "ps",
			Usage:   "pathwar [global flags] compose [compose flags] ps [flags]",
			FlagSet: composePSFlags,
			Options: []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(args []string) error {
				err := globalPreRun()
				if err != nil {
					return err
				}

				ctx := context.Background()
				cli, err := client.NewEnvClient()
				if err != nil {
					return errcode.ErrInitDockerClient.Wrap(err)
				}

				return pwcompose.PS(ctx, composePSDepth, cli, logger)
			},
		}, {
			Name:    "down",
			Usage:   "pathwar [global flags] compose [compose flags] down [flags] ID [ID...]",
			FlagSet: composeDownFlags,
			Options: []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(args []string) error {
				if err := globalPreRun(); err != nil {
					return err
				}

				ctx := context.Background()
				cli, err := client.NewEnvClient()
				if err != nil {
					return errcode.ErrInitDockerClient.Wrap(err)
				}

				composeCleanOpts.Logger = logger
				composeCleanOpts.ContainerIDs = args
				return pwcompose.Clean(ctx, cli, composeCleanOpts)
			},
		}},
		ShortHelp: "manage a challenge",
		FlagSet:   composeFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:      func([]string) error { return flag.ErrHelp },
	}
}
