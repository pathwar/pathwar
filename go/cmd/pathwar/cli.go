package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func cliCommand() *ffcli.Command {
	cliFlags := flag.NewFlagSet("cli", flag.ExitOnError)
	cliFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	cliFlags.StringVar(&ssoOpts.ClientID, "sso-client-id", ssoOpts.ClientID, "SSO ClientID")
	cliFlags.StringVar(&ssoOpts.ClientSecret, "sso-client-secret", ssoOpts.ClientSecret, "SSO ClientSecret")
	cliFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")
	cliFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")

	return &ffcli.Command{
		Name:      "cli",
		Usage:     "pathwar [global flags] cli [cli flags] <cmd> [cmd flags]",
		ShortHelp: "CLI replacement for the web portal",
		FlagSet:   cliFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:      func(args []string) error { return flag.ErrHelp },
		Subcommands: []*ffcli.Command{
			{
				Name:      "@me",
				ShortHelp: "Get an overview of your account (good place to start)",
				Exec: func(args []string) error {
					if err := globalPreRun(); err != nil {
						return err
					}
					ctx := context.Background()
					client, err := httpClientFromEnv(ctx)
					if err != nil {
						return err
					}
					var ret pwapi.UserGetSession_Output
					if err := client.RawProto(ctx, "GET", "/user/session", nil, &ret); err != nil {
						return err
					}

					logger.Debug("GET /user/session", zap.Any("ret", ret))

					// DB User
					{
						// fmt.Println(godev.PrettyJSON(ret.User))
						fmt.Printf("Welcome %s! 👋\n", ret.User.Username)
						createdAgo := humanize.Time(*ret.User.CreatedAt)
						updatedAgo := humanize.Time(*ret.User.UpdatedAt)
						fmt.Printf("Your account was created %s and updated %s.\n\n", createdAgo, updatedAgo)
					}

					// JWT Token
					{
						//fmt.Println(godev.PrettyJSON(ret.Claims))
						tokenAgo := humanize.Time(*ret.Claims.ActionToken.AuthTime)
						issuedAgo := humanize.Time(*ret.Claims.ActionToken.Iat)
						expireIn := humanize.Time(*ret.Claims.ActionToken.Exp)
						fmt.Printf("Your OAuth token was created %s, (re)issued %s and will expire in %s.\n\n", tokenAgo, issuedAgo, expireIn)
					}

					// Notifications
					{
						// FIXME: Todo
						// fmt.Println(godev.PrettyJSON(ret.Notifications))
					}

					// sessions
					{
						//fmt.Println(godev.PrettyJSON(ret.Seasons))
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"SEASON", "STATUS", "TEAM"})
						table.SetAlignment(tablewriter.ALIGN_CENTER)
						table.SetBorder(false)
						for _, entry := range ret.Seasons {
							name := entry.Season.Name
							// FIXME: use slug
							if entry.Season.IsDefault {
								name += " (default)"
							}
							status := ""
							switch entry.Season.Status {
							case pwdb.Season_Started:
								status = "Started 🏁"
							default:
								status = entry.Season.Status.String()
							}
							status += " / "
							switch entry.Season.Visibility {
							case pwdb.Season_Public:
								status += "Public 👐"
							default:
								status += entry.Season.Visibility.String()
							}
							team := "no team 👎"
							if entry.Team != nil {
								team = entry.Team.Organization.Name
							}
							table.Append([]string{name, status, team})
						}
						table.Render()
					}
					return nil
				},
			},
		},
	}
}
