package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"go.uber.org/zap"
	"moul.io/godev"
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

					var session pwapi.UserGetSession_Output
					if err := client.RawProto(ctx, "GET", "/user/session", nil, &session); err != nil {
						return err
					}
					logger.Debug("GET /user/session", zap.Any("ret", session))

					// DB User
					{
						// fmt.Println(godev.PrettyJSON(session.User))
						fmt.Printf("Welcome %s! 👋\n", session.User.Username)
						createdAgo := humanize.Time(*session.User.CreatedAt)
						updatedAgo := humanize.Time(*session.User.UpdatedAt)
						fmt.Printf("Your account was created %s and updated %s.\n\n", createdAgo, updatedAgo)
					}

					// JWT Token
					{
						//fmt.Println(godev.PrettyJSON(session.Claims))
						tokenAgo := humanize.Time(*session.Claims.ActionToken.AuthTime)
						issuedAgo := humanize.Time(*session.Claims.ActionToken.Iat)
						expireIn := humanize.Time(*session.Claims.ActionToken.Exp)
						fmt.Printf("Your OAuth token was created %s, (re)issued %s and will expire in %s.\n\n", tokenAgo, issuedAgo, expireIn)
					}

					// Notifications
					{
						// FIXME: Todo
						// fmt.Println(godev.PrettyJSON(session.Notifications))
					}

					// sessions
					{
						//fmt.Println(godev.PrettyJSON(session.Seasons))
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"SEASON", "STATUS", "TEAM"})
						table.SetAlignment(tablewriter.ALIGN_CENTER)
						table.SetBorder(false)
						for _, entry := range session.Seasons {
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
			}, {
				Name:      "teams",
				ShortHelp: "List teams, scores, etc",
				Exec: func(args []string) error {
					if err := globalPreRun(); err != nil {
						return err
					}
					ctx := context.Background()
					client, err := httpClientFromEnv(ctx)
					if err != nil {
						return err
					}
					var session pwapi.UserGetSession_Output
					if err := client.RawProto(ctx, "GET", "/user/session", nil, &session); err != nil {
						return err
					}
					logger.Debug("GET /user/session", zap.Any("ret", session))

					for _, seasonEntry := range session.Seasons {
						fmt.Printf("Season: %s\n", seasonEntry.Season.Name)
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"TEAM", "SCORE", "MEDALS", "ACHIEVEMENTS"})
						table.SetAlignment(tablewriter.ALIGN_CENTER)
						table.SetBorder(false)

						url := fmt.Sprintf("/teams?season_id=%d", seasonEntry.Season.ID)
						var ret pwapi.TeamList_Output
						if err := client.RawProto(ctx, "GET", url, nil, &ret); err != nil {
							return err
						}
						logger.Debug("GET "+url, zap.Any("ret", ret))

						for _, teamEntry := range ret.Items {
							name := teamEntry.Organization.Name
							score := fmt.Sprintf("%d", teamEntry.Score)
							achievements := fmt.Sprintf("%d", teamEntry.NbAchievements)
							medalParts := []string{}
							if nb := teamEntry.GoldMedals; nb > 0 {
								medalParts = append(medalParts, fmt.Sprintf("%d🥇", nb))
							}
							if nb := teamEntry.SilverMedals; nb > 0 {
								medalParts = append(medalParts, fmt.Sprintf("%d🥈", nb))
							}
							if nb := teamEntry.BronzeMedals; nb > 0 {
								medalParts = append(medalParts, fmt.Sprintf("%d🥉", nb))
							}
							medals := strings.Join(medalParts, " ")
							table.Append([]string{name, score, medals, achievements})
						}
						table.Render()
					}
					return nil
				},
			}, {
				Name:      "challenges",
				ShortHelp: "List challenges",
				Exec: func(args []string) error {
					if err := globalPreRun(); err != nil {
						return err
					}
					ctx := context.Background()
					client, err := httpClientFromEnv(ctx)
					if err != nil {
						return err
					}
					var session pwapi.UserGetSession_Output
					if err := client.RawProto(ctx, "GET", "/user/session", nil, &session); err != nil {
						return err
					}
					logger.Debug("GET /user/session", zap.Any("ret", session))

					for _, seasonEntry := range session.Seasons {
						fmt.Printf("Season: %s\n", seasonEntry.Season.Name)
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"CHALLENGE", "FLAVOR", "INSTANCE", "STATUS"})
						table.SetAlignment(tablewriter.ALIGN_CENTER)
						table.SetBorder(false)

						url := fmt.Sprintf("/season-challenges?season_id=%d", seasonEntry.Season.ID)
						var ret pwapi.SeasonChallengeList_Output
						if err := client.RawProto(ctx, "GET", url, nil, &ret); err != nil {
							return err
						}
						logger.Debug("GET "+url, zap.Any("ret", ret))

						for _, challengeEntry := range ret.Items {
							fmt.Println(godev.PrettyJSON(challengeEntry))
							name := challengeEntry.Flavor.Challenge.Name
							flavor := challengeEntry.Flavor.Version
							if challengeEntry.Flavor.IsLatest {
								flavor += " (latest)"
							}
							for _, instanceEntry := range challengeEntry.Flavor.Instances {
								instance := instanceEntry.NginxURL
								status := instanceEntry.Status.String()
								table.Append([]string{name, flavor, instance, status})
								name = ""
								flavor = ""
							}
						}
						table.Render()
					}
					return nil
				},
			},
		},
	}
}
