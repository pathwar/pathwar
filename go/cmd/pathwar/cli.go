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

var jsonFormat bool

func cliCommand() *ffcli.Command {
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
			cliMeCommand(),
			cliSeasonsCommand(),
			cliTeamsCommand(),
			cliChallengesCommand(),
			cliChallengeBuyCommand(),
			cliChallengeValidateCommand(),
			cliCouponValidateCommand(),
		},
	}
}

func cliMeCommand() *ffcli.Command {
	return &ffcli.Command{
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
	}
}

func cliSeasonsCommand() *ffcli.Command {
	return &ffcli.Command{
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
				table.SetHeader([]string{"SEASON", "NAME", "STATUS", "VISIBILITY", "MY TEAM", "CREATED", "UPDATED"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)
				for _, entry := range session.Seasons {
					slug := entry.Season.Slug
					name := entry.Season.Name
					// FIXME: use slug
					if entry.Season.IsGlobal {
						name += " 🌍"
					}
					status := ""
					switch entry.Season.Status {
					case pwdb.Season_Started:
						status = "Started 🏁"
					default:
						status = entry.Season.Status.String()
					}
					visibility := entry.Season.Visibility.String()
					if entry.Season.Visibility == pwdb.Season_Public {
						visibility = "Public 👐"
					}
					team := "no team 👎"
					if entry.Team != nil {
						team = entry.Team.Organization.Name
					}
					createdAgo := humanize.Time(*session.User.CreatedAt)
					updatedAgo := humanize.Time(*session.User.UpdatedAt)
					table.Append([]string{slug, name, status, visibility, team, createdAgo, updatedAgo})
				}
				table.Render()
			}
			return nil
		},
	}
}

func cliTeamsCommand() *ffcli.Command {
	return &ffcli.Command{
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
				table.SetHeader([]string{"TEAM", "NAME", "SCORE", "CASH", "MEDALS", "ACHIEVEMENTS"})
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
					slug := teamEntry.Organization.Slug
					name := teamEntry.Organization.Name
					score := fmt.Sprintf("%d", teamEntry.Score)
					cash := fmt.Sprintf("$%d", teamEntry.Cash)
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
					table.Append([]string{slug, name, score, cash, medals, achievements})
				}
				table.Render()
			}
			return nil
		},
	}
}

func cliChallengesCommand() *ffcli.Command {
	return &ffcli.Command{
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
				table.SetHeader([]string{"FLAVOR", "INSTANCES", "PRICE/REWARD", "SUBSCRIPTION", "URLS"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)
				table.SetColWidth(100)

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
					flavor := challengeEntry.Flavor
					flavorID := flavor.ASCIIID()
					subscription := "-"
					if len(challengeEntry.Subscriptions) > 0 {
						subscription = challengeEntry.Subscriptions[0].Status.String()
					}
					subscription = asciiStatus(subscription)
					instances := fmt.Sprintf("%d", len(flavor.Instances))
					urlParts := []string{}
					for _, instance := range flavor.Instances {
						if instance.NginxURL != "" {
							url := instance.NginxURL
							switch instance.Status {
							case pwdb.ChallengeInstance_Available:
								url += " 🟢"
							default:
								url += " 🔴"
							}
							urlParts = append(urlParts, url)
						}
					}
					urls := strings.Join(urlParts, ", ")
					price := "free"
					if flavor.PurchasePrice > 0 {
						price = fmt.Sprintf("$%d", flavor.PurchasePrice)
					}
					priceReward := fmt.Sprintf("%s / $%d", price, flavor.ValidationReward)
					table.Append([]string{flavorID, instances, priceReward, subscription, urls})
				}
				table.Render()
			}
			return nil
		},
	}
}

func cliChallengeBuyCommand() *ffcli.Command {
	input := pwapi.SeasonChallengeBuy_Input{}
	input.ApplyDefaults()
	flags := flag.NewFlagSet("cli challenge buy", flag.ExitOnError)
	flags.StringVar(&input.FlavorID, "flavor", input.FlavorID, "Flavor ID or Slug")
	flags.StringVar(&input.SeasonID, "season", input.SeasonID, "Season ID or Slug")

	return &ffcli.Command{
		Name:      "challenge-buy",
		Usage:     "pathwar [global flags] cli [cli flags] challenge-buy --flavor=XXX",
		FlagSet:   flags,
		ShortHelp: "Buy a challenge",
		Exec: func(args []string) error {
			if input.FlavorID == "" {
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

			var ret pwapi.SeasonChallengeBuy_Output
			err = client.RawProto(ctx, "POST", "/season-challenge/buy", &input, &ret)

			if jsonFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			logger.Debug("POST /season-challenge/buy", zap.Any("input", input), zap.Any("ret", ret), zap.Error(err))
			switch {
			case err == nil:
				fmt.Printf("%s: successfully bought (%d)\n", input.FlavorID, ret.ChallengeSubscription.ID)
			case strings.Contains(err.Error(), "ErrChallengeAlreadySubscribed(#4011)"):
				fmt.Printf("%s: already bought\n", input.FlavorID)
			default:
				return err
			}
			return nil
		},
	}
}

func cliChallengeValidateCommand() *ffcli.Command {
	input := pwapi.ChallengeSubscriptionValidate_Input{}

	passphrases := ""
	flags := flag.NewFlagSet("cli challenge validate", flag.ExitOnError)
	flags.StringVar(&passphrases, "passphrases", passphrases, "Passphrases separated with commas")
	flags.Int64Var(&input.ChallengeSubscriptionID, "subscription-id", input.ChallengeSubscriptionID, "Challenge subscription ID")
	flags.StringVar(&input.Comment, "comment", input.Comment, "Comment for validation")

	return &ffcli.Command{
		Name:      "challenge-validate",
		Usage:     "pathwar [global flags] cli [cli flags] challenge-validate [flags]",
		FlagSet:   flags,
		ShortHelp: "Validate a challenge",
		Exec: func(args []string) error {
			if input.ChallengeSubscriptionID == 0 {
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

			input.Passphrases = strings.Split(passphrases, ",")
			var ret pwapi.ChallengeSubscriptionValidate_Output
			err = client.RawProto(ctx, "POST", "/challenge-subscription/validate", &input, &ret)
			if err != nil {
				return err
			}
			logger.Debug("POST /challenge-subscription/validate", zap.Any("input", input), zap.Any("ret", ret), zap.Error(err))
			fmt.Println(godev.PrettyJSONPB(&ret))
			return nil
		},
	}
}

func cliCouponValidateCommand() *ffcli.Command {
	return &ffcli.Command{
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
	}
}
