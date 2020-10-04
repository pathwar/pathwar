package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/docker/docker/client"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"gopkg.in/yaml.v2"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
)

func composeCommand() *ffcli.Command {
	var composeFlags = flag.NewFlagSet("compose", flag.ExitOnError)
	return &ffcli.Command{
		Name:  "compose",
		Usage: "pathwar [global flags] compose [compose flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{
			composeUpCommand(),
			composePrepareCommand(),
			composePsCommand(),
			composeDownCommand(),
			composeRegisterCommand(),
		},
		ShortHelp: "manage a challenge",
		FlagSet:   composeFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:      func([]string) error { return flag.ErrHelp },
	}
}

func composeUpCommand() *ffcli.Command {
	var (
		composeUpOpts  = pwcompose.NewUpOpts()
		composeUpFlags = flag.NewFlagSet("compose up", flag.ExitOnError)
	)
	composeUpFlags.StringVar(&composeUpOpts.InstanceKey, "instance-key", composeUpOpts.InstanceKey, "instance key used to generate instance ID")
	composeUpFlags.BoolVar(&composeUpOpts.ForceRecreate, "force-recreate", composeUpOpts.ForceRecreate, "down previously created instances of challenge")
	return &ffcli.Command{
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
	}
}

func composePrepareCommand() *ffcli.Command {
	var (
		composePrepareOpts  = pwcompose.NewPrepareOpts()
		composePrepareFlags = flag.NewFlagSet("compose prepare", flag.ExitOnError)
	)
	composePrepareFlags.BoolVar(&composePrepareOpts.NoPush, "no-push", composePrepareOpts.NoPush, "don't push images")
	composePrepareFlags.StringVar(&composePrepareOpts.Prefix, "prefix", composePrepareOpts.Prefix, "docker image prefix")
	composePrepareFlags.StringVar(&composePrepareOpts.Version, "version", composePrepareOpts.Version, "challenge version")
	composePrepareFlags.BoolVar(&composePrepareOpts.JSON, "json", composePrepareOpts.JSON, "JSON format")
	return &ffcli.Command{
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
	}
}

func composePsCommand() *ffcli.Command {
	var composePSFlags = flag.NewFlagSet("compose ps", flag.ExitOnError)
	composePSFlags.IntVar(&composePSDepth, "depth", 0, "depth to display")
	return &ffcli.Command{
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
	}
}

func composeDownCommand() *ffcli.Command {
	var (
		composeDownFlags = flag.NewFlagSet("compose down", flag.ExitOnError)
		composeCleanOpts = pwcompose.NewCleanOpts()
	)
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveVolumes, "rm-volumes", composeCleanOpts.RemoveVolumes, "keep volumes")
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveImages, "rm-images", composeCleanOpts.RemoveImages, "remove images as well")
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveNginx, "rm-nginx", composeCleanOpts.RemoveNginx, "down nginx container and proxy network as well")
	return &ffcli.Command{
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
	}
}

func composeRegisterCommand() *ffcli.Command {
	var (
		registerPrint        bool
		composeRegisterFlags = flag.NewFlagSet("compose register", flag.ExitOnError)
	)
	composeRegisterFlags.BoolVar(&registerPrint, "print", false, "print pathwar commands (pipe-friendly)")
	return &ffcli.Command{
		Name:    "register",
		Usage:   "pathwar [global flags] compose [compose flags] register [flags] path/to/pathwar-compose.yml",
		FlagSet: composeRegisterFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec: func(args []string) error {
			if len(args) != 1 {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			if !registerPrint {
				return errors.New("--print flag is required (for now)")
			}

			composePath := args[0]
			f, err := os.Open(composePath)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			defer f.Close()

			content, err := ioutil.ReadAll(f)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			var config pwcompose.PathwarConfig
			if err = yaml.Unmarshal(content, &config); err != nil {
				return errcode.TODO.Wrap(err)
			}

			// fmt.Println(godev.PrettyJSON(config))

			slug := config.Pathwar.Challenge.Slug
			if slug == "" {
				return errors.New("a challenge slug is required in docker-compose.yml")
			}

			command := []string{"pathwar", "admin", "challenge-add"}
			if slug != "" {
				command = append(command, "--slug", shellescape.Quote(slug))
			}
			if name := config.Pathwar.Challenge.Name; name != "" {
				command = append(command, "--name", shellescape.Quote(name))
			}
			if description := config.Pathwar.Challenge.Description; description != "" {
				command = append(command, "--description", shellescape.Quote(description))
			}
			if homepage := config.Pathwar.Challenge.Homepage; homepage != "" {
				command = append(command, "--homepage", shellescape.Quote(homepage))
			}
			if locale := config.Pathwar.Challenge.Locale; locale != "" {
				command = append(command, "--locale", shellescape.Quote(locale))
			}
			if author := config.Pathwar.Challenge.Author; author != "" {
				command = append(command, "--author", shellescape.Quote(author))
			}
			fmt.Println(strings.Join(command, " "))

			command = []string{"pathwar", "admin", "challenge-flavor-add", "--challenge", slug}
			if slug := config.Pathwar.Flavor.Slug; slug != "" {
				command = append(command, "--slug", shellescape.Quote(slug))
			}
			if body := config.Pathwar.Flavor.Body; body != "" {
				command = append(command, "--body", shellescape.Quote(body))
			}
			if version := config.Pathwar.Flavor.Version; version != "" {
				command = append(command, "--version", shellescape.Quote(version))
			}
			if passphrases := config.Pathwar.Flavor.Passphrases; passphrases != 0 {
				command = append(command, "--passphrases", fmt.Sprintf("%d", passphrases))
			}
			if category := config.Pathwar.Flavor.Category; category != "" {
				command = append(command, "--category", shellescape.Quote(category))
			}
			if redumpPolicy := config.Pathwar.Flavor.RedumpPolicy; redumpPolicy != nil {
				policy, _ := json.Marshal(config.Pathwar.Flavor.RedumpPolicy)
				command = append(command, "--redump-policy", shellescape.Quote(string(policy)))
			}
			if tags := config.Pathwar.Flavor.Tags; len(tags) > 0 {
				command = append(command, "--tags", shellescape.Quote(strings.Join(tags, ",")))
			}
			if validationReward := config.Pathwar.Flavor.ValidationReward; validationReward != 0 {
				command = append(command, "--validation-reward", fmt.Sprintf("%d", validationReward))
			}
			if purchasePrice := config.Pathwar.Flavor.PurchasePrice; purchasePrice != 0 {
				command = append(command, "--purchase-price", fmt.Sprintf("%d", purchasePrice))
			}
			command = append(command, "--compose-bundle", shellescape.Quote(composePath))
			fmt.Println(strings.Join(command, " "))

			// FIXME: add a bool in the yaml file to make it dynamic
			command = []string{"pathwar", "admin", "season-challenge-add", "--flavor", slug, "--season", "global"}
			fmt.Println(strings.Join(command, " "))
			return nil
		},
	}
}
