package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"moul.io/godev"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
)

func adminCommand() *ffcli.Command {
	var (
		adminFlags                     = flag.NewFlagSet("admin", flag.ExitOnError)
		adminPSFlags                   = flag.NewFlagSet("admin ps", flag.ExitOnError)
		adminRedumpFlags               = flag.NewFlagSet("admin redump", flag.ExitOnError)
		adminChallengeAddFlags         = flag.NewFlagSet("admin challenge add", flag.ExitOnError)
		adminChallengeFlavorAddFlags   = flag.NewFlagSet("admin challenge flavor add", flag.ExitOnError)
		adminChallengeInstanceAddFlags = flag.NewFlagSet("admin challenge instance add", flag.ExitOnError)
		jsonFormat                     bool
	)
	adminFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	adminFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")
	adminFlags.BoolVar(&jsonFormat, "json", false, "Print JSON and exit")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Name, "name", "", "Challenge name")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Description, "description", "", "Challenge description")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Author, "author", "", "Challenge author")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Locale, "locale", "", "Challenge Locale")
	adminChallengeAddFlags.BoolVar(&adminChallengeAddInput.Challenge.IsDraft, "is-draft", true, "Is challenge production ready ?")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.PreviewUrl, "preview-url", "", "Challenge preview URL")
	adminChallengeAddFlags.StringVar(&adminChallengeAddInput.Challenge.Homepage, "homepage", "", "Challenge homepage URL")
	adminChallengeFlavorAddFlags.StringVar(&adminChallengeFlavorAddInput.ChallengeFlavor.Version, "version", "1.0.0", "Challenge flavor version")
	adminChallengeFlavorAddFlags.StringVar(&adminChallengeFlavorAddInput.ChallengeFlavor.ComposeBundle, "compose-bundle", "", "Challenge flavor compose bundle")
	adminChallengeFlavorAddFlags.Int64Var(&adminChallengeFlavorAddInput.ChallengeFlavor.ChallengeID, "challenge-id", 0, "Challenge id")
	adminChallengeInstanceAddFlags.Int64Var(&adminChallengeInstanceAddInput.ChallengeInstance.AgentID, "agent-id", 0, "Id of the agent that will host the instance")
	adminChallengeInstanceAddFlags.Int64Var(&adminChallengeInstanceAddInput.ChallengeInstance.FlavorID, "flavor-id", 0, "Challenge flavor id")

	return &ffcli.Command{
		Name:  "admin",
		Usage: "pathwar [global flags] admin [admin flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{{
			Name:    "ps",
			Usage:   "pathwar [global flags] admin [admin flags] ps [flags]",
			FlagSet: adminPSFlags,
			Exec: func(args []string) error {
				if err := globalPreRun(); err != nil {
					return err
				}

				ctx := context.Background()
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				ret, err := apiClient.AdminPS(ctx, &pwapi.AdminPS_Input{})
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				if jsonFormat {
					fmt.Println(godev.PrettyJSONPB(&ret))
					return nil
				}

				return nil
			},
		}, {
			Name:    "redump",
			Usage:   "pathwar [global flags] admin [admin flags] redump [flags] ID...",
			FlagSet: adminRedumpFlags,
			Exec: func(args []string) error {
				if len(args) < 1 {
					return flag.ErrHelp
				}

				if err := globalPreRun(); err != nil {
					return err
				}

				ctx := context.Background()
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				_, err = apiClient.AdminRedump(ctx, &pwapi.AdminRedump_Input{
					Identifiers: args,
				})
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				fmt.Println("OK")

				return nil
			},
		}, {
			Name:      "challenge-add",
			Usage:     "pathwar [global flags] admin [admin flags] challenge-add [flags] [args...]",
			ShortHelp: "add a challenge",
			FlagSet:   adminChallengeAddFlags,
			Exec: func(args []string) error {
				if err := globalPreRun(); err != nil {
					return err
				}

				ctx := context.Background()
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				ret, err := apiClient.AdminAddChallenge(ctx, &adminChallengeAddInput)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}
				if globalDebug {
					fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
				}
				fmt.Println(ret.Challenge.ID)
				return nil
			},
		}, {
			Name:      "challenge-flavor-add",
			Usage:     "pathwar [global flags] admin [admin flags] challenge-flavor-add [flags] [args...]",
			ShortHelp: "add a challenge flavor",
			FlagSet:   adminChallengeFlavorAddFlags,
			Exec: func(args []string) error {
				if err := globalPreRun(); err != nil {
					return err
				}

				ctx := context.Background()
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				ret, err := apiClient.AdminAddChallengeFlavor(ctx, &adminChallengeFlavorAddInput)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}
				if globalDebug {
					fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
				}
				fmt.Println(ret.ChallengeFlavor.ID)
				return nil
			},
		}, {
			Name:      "challenge-instance-add",
			Usage:     "pathwar [global flags] admin [admin flags] challenge-instance-add [flags] [args...]",
			ShortHelp: "add a challenge instance",
			FlagSet:   adminChallengeInstanceAddFlags,
			Exec: func(args []string) error {
				if err := globalPreRun(); err != nil {
					return err
				}

				ctx := context.Background()
				apiClient, err := httpClientFromEnv(ctx)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}

				ret, err := apiClient.AdminAddChallengeInstance(ctx, &adminChallengeInstanceAddInput)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}
				if globalDebug {
					fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
				}
				fmt.Println(ret.ChallengeInstance.ID)
				return nil
			},
		}},
		ShortHelp: "admin commands",
		FlagSet:   adminFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:      func([]string) error { return flag.ErrHelp },
	}
}
