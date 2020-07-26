package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
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
	var jsonFormat bool
	cliFlags := flag.NewFlagSet("cli", flag.ExitOnError)
	cliFlags.BoolVar(&jsonFormat, "json", false, "Print JSON output")
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

					if jsonFormat {
						fmt.Println(godev.PrettyJSONPB(&session))
						return nil
					}

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
						fmt.Printf("Your OAuth token was created %s, (re)issued %s and will expire in %s.\n", tokenAgo, issuedAgo, expireIn)
					}

					// Notifications
					{
						// FIXME: Todo
						// fmt.Println(godev.PrettyJSON(session.Notifications))
					}
					return nil
				},
			}, {
				Name:      "seasons",
				ShortHelp: "List available seasons",
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

					if jsonFormat {
						fmt.Println(godev.PrettyJSON(session.Seasons))
						return nil
					}

					// seasons
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
						if !jsonFormat {
							fmt.Printf("Season: %s\n", seasonEntry.Season.Name)
						}
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"TEAM", "SCORE", "MEDALS", "ACHIEVEMENTS"})
						table.SetAlignment(tablewriter.ALIGN_CENTER)
						table.SetBorder(false)

						url := fmt.Sprintf("/teams?season_id=%d", seasonEntry.Season.ID)
						var ret pwapi.TeamList_Output
						if err := client.RawProto(ctx, "GET", url, nil, &ret); err != nil {
							return err
						}

						if jsonFormat {
							fmt.Println(godev.PrettyJSONPB(&ret))
							continue
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
						if !jsonFormat {
							fmt.Printf("Season: %s\n", seasonEntry.Season.Name)
						}
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"ID", "CHALLENGE", "FLAVOR", "URL", "STATUS", "SUBSCRIPTION"})
						table.SetAlignment(tablewriter.ALIGN_CENTER)
						table.SetBorder(false)

						url := fmt.Sprintf("/season-challenges?season_id=%d", seasonEntry.Season.ID)
						var ret pwapi.SeasonChallengeList_Output
						if err := client.RawProto(ctx, "GET", url, nil, &ret); err != nil {
							return err
						}

						if jsonFormat {
							fmt.Println(godev.PrettyJSONPB(&ret))
							continue
						}

						logger.Debug("GET "+url, zap.Any("ret", ret))

						for _, challengeEntry := range ret.Items {
							id := fmt.Sprintf("%d", challengeEntry.ID)
							name := challengeEntry.Flavor.Challenge.Name
							flavor := challengeEntry.Flavor.Version
							if challengeEntry.Flavor.IsLatest {
								flavor += " (latest)"
							}
							subscription := "none"
							if len(challengeEntry.Subscriptions) > 0 {
								subscription = challengeEntry.Subscriptions[0].Status.String()
							}
							for _, instanceEntry := range challengeEntry.Flavor.Instances {
								instance := instanceEntry.NginxURL
								status := instanceEntry.Status.String()
								if instanceEntry.Status == pwdb.ChallengeInstance_Available {
									status += " 👍"
								}
								table.Append([]string{id, name, flavor, instance, status, subscription})
								// merge name and flavor if multiple subscription
								name = ""
								flavor = ""
								subscription = ""
							}
						}
						table.Render()
					}
					return nil
				},
			}, {
				Name:      "challenge-buy",
				Usage:     "pathwar [global flags] cli [cli flags] challenge-buy ID...",
				ShortHelp: "Buy a challenge",
				Exec: func(args []string) error {
					if len(args) < 1 {
						return flag.ErrHelp
					}

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

					for _, arg := range args {
						id, err := strconv.Atoi(arg)
						if err != nil {
							return err
						}

						input := pwapi.SeasonChallengeBuy_Input{
							SeasonChallengeID: int64(id),
							TeamID:            session.User.ActiveTeamMember.TeamID, // FIXME: dynamic
						}
						var ret pwapi.SeasonChallengeBuy_Output
						err = client.RawProto(ctx, "POST", "/season-challenge/buy", &input, &ret)

						if jsonFormat {
							fmt.Println(godev.PrettyJSONPB(&ret))
							continue
						}

						logger.Debug("POST /season-challenge/buy", zap.Any("input", input), zap.Any("ret", ret), zap.Error(err))
						switch {
						case err == nil:
							fmt.Printf("%d: successfully bought\n", id)
						case strings.Contains(err.Error(), "ErrChallengeAlreadySubscribed(#4011)"):
							fmt.Printf("%d: already bought\n", id)
						default:
							return err
						}
					}
					return nil
				},
			}, {
				Name:      "coupon-validate",
				Usage:     "pathwar [global flags] cli [cli flags] coupon-validate COUPON...",
				ShortHelp: "Validate a coupon",
				Exec: func(args []string) error {
					if len(args) < 1 {
						return flag.ErrHelp
					}

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

					for _, arg := range args {
						input := pwapi.CouponValidate_Input{
							Hash:   arg,
							TeamID: session.User.ActiveTeamMember.TeamID, // FIXME: dynamic
						}
						var ret pwapi.CouponValidate_Output
						err = client.RawProto(ctx, "POST", "/coupon-validation", &input, &ret)

						if jsonFormat {
							fmt.Println(godev.PrettyJSONPB(&ret))
							continue
						}

						logger.Debug("POST /coupon-validation", zap.Any("input", input), zap.Any("ret", ret), zap.Error(err))
						switch {
						case err == nil:
							fmt.Printf("coupon %q validated\n", arg)
						case strings.Contains(err.Error(), "ErrCouponNotFound(#4063)"):
							fmt.Printf("coupon %q does not exist\n", arg)
						default:
							return err
						}
					}
					return nil
				},
			},
		},
	}
}
