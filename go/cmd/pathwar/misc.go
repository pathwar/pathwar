package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

func miscCommand() *ffcli.Command {
	ssoFlags := flag.NewFlagSet("sso", flag.ExitOnError)
	ssoFlags.BoolVar(&ssoOpts.AllowUnsafe, "unsafe", ssoOpts.AllowUnsafe, "Allow unsafe SSO")
	ssoFlags.StringVar(&ssoOpts.ClientID, "clientid", ssoOpts.ClientID, "SSO ClientID")
	ssoFlags.StringVar(&ssoOpts.Pubkey, "pubkey", ssoOpts.Pubkey, "SSO Public Key")
	ssoFlags.StringVar(&ssoOpts.Realm, "realm", ssoOpts.Realm, "SSO Realm")

	return &ffcli.Command{
		Name:       "misc",
		ShortUsage: "pathwar [global flags] misc [misc flags] <subcommand> [flags] [args...]",
		ShortHelp:  "misc contains advanced commands",
		Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:       func(ctx context.Context, args []string) error { return flag.ErrHelp },
		Subcommands: []*ffcli.Command{
			{
				Name:       "pwinit-binary",
				ShortUsage: "pathwar [global flags] misc [misc flags] pwinit-binary",
				Exec: func(ctx context.Context, args []string) error {
					binary, err := pwinit.Binary()
					if err != nil {
						return err
					}
					os.Stdout.Write(binary)
					return nil
				},
			}, {
				Name:       "sso",
				ShortUsage: "pathwar [global flags] sso [sso flags] <subcommand> [flags] [args...]",
				ShortHelp:  "manage SSO tokens",
				FlagSet:    ssoFlags,
				Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
				Exec:       func(ctx context.Context, args []string) error { return flag.ErrHelp },
				Subcommands: []*ffcli.Command{
					{
						Name:       "token",
						ShortUsage: "pathwar [global flags] sso [sso flags] token TOKEN",
						Exec: func(ctx context.Context, args []string) error {
							if len(args) < 1 {
								return flag.ErrHelp
							}
							if err := globalPreRun(); err != nil {
								return err
							}
							sso, err := ssoFromFlags()
							if err != nil {
								return errcode.ErrGetSSOClientFromFlags.Wrap(err)
							}

							// token
							token, _, err := sso.TokenWithClaims(args[0])
							if err != nil {
								return errcode.ErrGetSSOClaims.Wrap(err)
							}
							out, _ := json.MarshalIndent(token, "", "  ")
							fmt.Println(string(out))

							return nil
						},
					}, {
						Name:       "logout",
						ShortUsage: "pathwar [global flags] sso [sso flags] logout TOKEN",
						Exec: func(ctx context.Context, args []string) error {
							if len(args) < 1 {
								return flag.ErrHelp
							}
							err := globalPreRun()
							if err != nil {
								return err
							}

							sso, err := ssoFromFlags()
							if err != nil {
								return errcode.ErrGetSSOClientFromFlags.Wrap(err)
							}

							// logout
							if err := sso.Logout(args[0]); err != nil {
								return errcode.ErrGetSSOLogout.Wrap(err)
							}
							return nil
						},
					}, {
						Name:       "whoami",
						ShortUsage: "pathwar [global flags] sso [sso flags] whoami TOKEN",
						Exec: func(ctx context.Context, args []string) error {
							if len(args) < 1 {
								return flag.ErrHelp
							}
							err := globalPreRun()
							if err != nil {
								return err
							}

							sso, err := ssoFromFlags()
							if err != nil {
								return errcode.ErrGetSSOClientFromFlags.Wrap(err)
							}

							// whoami
							info, err := sso.Whoami(args[0])
							if err != nil {
								return errcode.ErrGetSSOWhoami.Wrap(err)
							}
							for k, v := range info {
								fmt.Printf("- %s: %v\n", k, v)
							}
							return nil
						},
					},
				},
			},
		},
	}
}
