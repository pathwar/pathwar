package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"moul.io/godev"
	"moul.io/zapconfig"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

var adminJSONFormat bool

func adminCommand() *ffcli.Command {
	adminFlags := flag.NewFlagSet("admin", flag.ExitOnError)
	adminFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	adminFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")
	adminFlags.BoolVar(&adminJSONFormat, "json", false, "Print JSON and exit")
	return &ffcli.Command{
		Name:       "admin",
		ShortUsage: "pathwar [global flags] admin [admin flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{
			// read-only
			adminChallengesCommand(),
			adminUsersCommand(),
			adminAgentsCommand(),
			adminActivitiesCommand(),
			adminOrganizationsCommand(),
			adminTeamsCommand(),
			adminCouponsCommand(),
			adminChallengeSubscriptionsCommand(),
			adminListAllCommand(),
			adminSearchCommand(),

			// actions
			adminAddCouponCommand(),
			adminRedumpCommand(),
			adminChallengeAddCommand(),
			adminChallengeRedumpCommand(),
			adminChallengeFlavorAddCommand(),
			adminSeasonAddCommand(),
			adminSeasonChallengeAddCommand(),
			adminTestingSeasonUserAddCommand(),
		},
		ShortHelp: "admin commands",
		FlagSet:   adminFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:      func(ctx context.Context, args []string) error { return flag.ErrHelp },
	}
}

func adminChallengesCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin challenges", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "challenges",
		ShortUsage: "pathwar [global flags] admin [admin flags] challenges [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListChallenges(ctx, &pwapi.AdminListChallenges_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// tree
			{
				fmt.Println("TREE")
				for _, challenge := range ret.Challenges {
					fmt.Printf("- %s\n", challenge.Slug)
					for _, flavor := range challenge.Flavors {
						fmt.Printf("  - %s\n", flavor.Slug)
						for _, instance := range flavor.Instances {
							fmt.Printf("    - %d\n", instance.ID)
						}
					}
				}
				fmt.Println("")
			}

			// challenges table
			{
				fmt.Println("CHALLENGES")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"CHALLENGE", "NAME", "AUTHOR", "CREATED", "UPDATED", "FLAVORS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, challenge := range ret.Challenges {
					slug := challenge.Slug
					name := challenge.Name
					author := challenge.Author
					createdAgo := humanize.Time(*challenge.CreatedAt)
					updatedAgo := humanize.Time(*challenge.UpdatedAt)
					flavors := fmt.Sprintf("%d", len(challenge.Flavors))
					id := fmt.Sprintf("%d", challenge.ID)
					table.Append([]string{slug, name, author, createdAgo, updatedAgo, flavors, id})
				}
				table.Render()
				fmt.Println("")
			}

			// flavors table
			{
				fmt.Println("FLAVORS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"FLAVOR", "CREATED", "UPDATED", "INSTANCES", "SEASON CHALLENGES", "PRICE/REWARD", "PASSPHRASES", "BODY"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, challenge := range ret.Challenges {
					for _, flavor := range challenge.Flavors {
						slug := flavor.Slug
						createdAgo := humanize.Time(*flavor.CreatedAt)
						updatedAgo := humanize.Time(*flavor.UpdatedAt)
						instances := asciiInstancesStats(flavor.Instances)
						seasonChallengeParts := []string{}
						for _, seasonChallenge := range flavor.SeasonChallenges {
							seasonChallengeParts = append(seasonChallengeParts, seasonChallenge.Season.Slug)
						}
						seasonChallenges := strings.Join(seasonChallengeParts, ", ")
						price := "free"
						if flavor.PurchasePrice > 0 {
							price = fmt.Sprintf("$%d", flavor.PurchasePrice)
						}
						priceReward := fmt.Sprintf("%s / $%d", price, flavor.ValidationReward)
						body := flavor.Body
						if len(body) > 16 {
							body = body[:15] + "..."
						}
						passphrases := fmt.Sprintf("%d", flavor.Passphrases)
						slug = strings.TrimSuffix(slug, "@default")
						table.Append([]string{slug, createdAgo, updatedAgo, instances, seasonChallenges, priceReward, passphrases, body})
					}
				}
				table.Render()
				fmt.Println("")
			}

			// instances
			{
				fmt.Println("INSTANCES")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"INSTANCE", "FLAVOR", "AGENT", "STATUS", "CREATED", "UPDATED", "CONFIG", "SEASON CHALLENGES"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, challenge := range ret.Challenges {
					for _, flavor := range challenge.Flavors {
						for _, instance := range flavor.Instances {
							// fmt.Println(godev.PrettyJSONPB(instance))
							id := fmt.Sprintf("%d", instance.ID)
							status := asciiStatus(instance.Status.String())
							agentSlug := instance.Agent.ASCIIID()
							flavorSlug := flavor.Slug
							createdAgo := humanize.Time(*instance.CreatedAt)
							updatedAgo := humanize.Time(*instance.UpdatedAt)
							configStruct, _ := instance.ParseInstanceConfig()
							config := godev.JSONPB(configStruct)
							if len(config) > 30 {
								config = config[:28] + "..."
							}

							seasonChallenges := fmt.Sprintf("%d", len(flavor.SeasonChallenges))
							table.Append([]string{id, flavorSlug, agentSlug, status, createdAgo, updatedAgo, config, seasonChallenges})
						}
					}
				}
				table.Render()
			}

			return nil
		},
	}
}

func adminUsersCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin users", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "users",
		ShortUsage: "pathwar [global flags] admin [admin flags] users [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListUsers(ctx, &pwapi.AdminListUsers_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// users table
			{
				fmt.Println("USERS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"USER", "STATUS", "USERNAME", "EMAIL", "CREATED", "UPDATED", "TEAMS", "ORGS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, user := range ret.Users {
					// fmt.Println(godev.PrettyJSONPB(user))
					slug := user.Slug
					email := user.Email
					username := user.Username
					status := asciiStatus(user.DeletionStatus.String())
					createdAgo := humanize.Time(*user.CreatedAt)
					updatedAgo := humanize.Time(*user.UpdatedAt)
					teams := fmt.Sprintf("%d", len(user.TeamMemberships))
					orgs := fmt.Sprintf("%d", len(user.OrganizationMemberships))
					id := fmt.Sprintf("%d", user.ID)
					table.Append([]string{slug, status, username, email, createdAgo, updatedAgo, teams, orgs, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminAgentsCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin agents", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "agents",
		ShortUsage: "pathwar [global flags] admin [admin flags] agents [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListAgents(ctx, &pwapi.AdminListAgents_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// agents table
			{
				fmt.Println("AGENTS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"AGENT", "HOSTNAME", "SUFFIX", "STATUS", "CREATED", "UPDATED", "SEEN", "STATS", "INSTANCES", "DEFAULT", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, agent := range ret.Agents {
					// fmt.Println(godev.PrettyJSONPB(agent))
					slug := agent.Slug
					createdAgo := humanize.Time(*agent.CreatedAt)
					updatedAgo := humanize.Time(*agent.UpdatedAt)
					seenAgo := humanize.Time(*agent.LastSeenAt)
					id := fmt.Sprintf("%d", agent.ID)
					instances := asciiInstancesStats(agent.ChallengeInstances)
					stats := fmt.Sprintf("%d seen / %d reg.", agent.TimesSeen, agent.TimesRegistered)
					status := asciiStatus(agent.Status.String())
					isDefault := asciiBool(agent.DefaultAgent)
					suffix := agent.DomainSuffix
					hostname := agent.Hostname
					table.Append([]string{slug, hostname, suffix, status, createdAgo, updatedAgo, seenAgo, stats, instances, isDefault, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminActivitiesCommand() *ffcli.Command {
	input := pwapi.AdminListActivities_Input{
		FilteringPreset: "default",
		Limit:           50,
		Since:           nil,
	}
	flags := flag.NewFlagSet("admin activities", flag.ExitOnError)
	var (
		follow    bool
		logFormat bool
	)
	flags.BoolVar(&follow, "f", follow, "follow events")
	flags.BoolVar(&logFormat, "log", logFormat, "use a logging format")
	flags.Int64Var(&input.Limit, "limit", input.Limit, "max entries per call")
	flags.StringVar(&input.FilteringPreset, "preset", input.FilteringPreset, "filtering preset (registers)")

	return &ffcli.Command{
		Name:       "activities",
		ShortUsage: "pathwar [global flags] admin [admin flags] activities [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListActivities(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			// adminJSONFormat is ignored if -log is provided
			if !logFormat && adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			if logFormat {
				logger := zapconfig.New().SetPreset("light-console").MustBuild().WithOptions(zap.WithCaller(false))
				for _, activity := range ret.Activities {
					activity.Log(logger)
				}
			} else {
				fmt.Println("ACTIVITIES")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"ID", "KIND", "HAPPENED", "AUTHOR", "TEAM", "USER", "AGENT", "ORG", "SEASON", "CHALLENGE", "FLAVOR", "INSTANCE", "COUPON", "SEASON CHAL.", "TEAM MEMBER", "CHALLENGE SUBS"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, activity := range ret.Activities {
					// fmt.Println(godev.PrettyJSONPB(activity))
					author := activity.Author.Slug
					id := fmt.Sprintf("%d", activity.ID)
					kind := activity.Kind.String()
					createdAgo := humanize.Time(*activity.CreatedAt)
					team := activity.Team.ASCIIID()
					user := activity.User.ASCIIID()
					agent := activity.Agent.ASCIIID()
					organization := activity.Organization.ASCIIID()
					season := activity.Season.ASCIIID()
					challenge := activity.Challenge.ASCIIID()
					flavor := activity.ChallengeFlavor.ASCIIID()
					instance := activity.ChallengeInstance.ASCIIID()
					coupon := activity.Coupon.ASCIIID()
					seasonChallenge := activity.SeasonChallenge.ASCIIID()
					teamMember := activity.TeamMember.ASCIIID()
					challengeSubscription := activity.ChallengeSubscription.ASCIIID()
					table.Append([]string{id, kind, createdAgo, author, team, user, agent, organization, season, challenge, flavor, instance, coupon, seasonChallenge, teamMember, challengeSubscription})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminOrganizationsCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin organizations", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "organizations",
		ShortUsage: "pathwar [global flags] admin [admin flags] organizations [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListOrganizations(ctx, &pwapi.AdminListOrganizations_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// organizations table
			{
				fmt.Println("ORGANIZATIONS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"SLUG", "NAME", "STATUS", "CREATED", "UPDATED", "TEAMS", "MEMBERS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, organization := range ret.Organizations {
					// fmt.Println(godev.PrettyJSONPB(organization))
					id := fmt.Sprintf("%d", organization.ID)
					createdAgo := humanize.Time(*organization.CreatedAt)
					updatedAgo := humanize.Time(*organization.UpdatedAt)
					slug := organization.Slug
					name := organization.Name
					status := asciiStatus(organization.DeletionStatus.String())
					teamParts := []string{}
					for _, team := range organization.Teams {
						teamParts = append(teamParts, team.Season.ASCIIID())
					}
					teams := strings.Join(teamParts, ",")
					memberParts := []string{}
					for _, member := range organization.Members {
						memberParts = append(memberParts, member.User.ASCIIID())
					}
					members := strings.Join(memberParts, ",")
					table.Append([]string{slug, name, status, createdAgo, updatedAgo, teams, members, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminTeamsCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin teams", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "teams",
		ShortUsage: "pathwar [global flags] admin [admin flags] teams [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListTeams(ctx, &pwapi.AdminListTeams_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// teams table
			{
				fmt.Println("TEAMS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"SLUG", "CREATED", "UPDATED", "SEASON", "CASH", "MEMBERS", "ORGANIZATION", "CHALLENGES", "ACHIEVEMENTS", "STATUS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, team := range ret.Teams {
					// fmt.Println(godev.PrettyJSONPB(team))
					id := fmt.Sprintf("%d", team.ID)
					createdAgo := humanize.Time(*team.CreatedAt)
					updatedAgo := humanize.Time(*team.UpdatedAt)
					season := team.Season.ASCIIID()
					memberParts := []string{}
					for _, member := range team.Members {
						memberParts = append(memberParts, member.User.ASCIIID())
					}
					members := strings.Join(memberParts, ",")
					organization := team.Organization.ASCIIID()
					status := asciiStatus(team.DeletionStatus.String())
					challenges := asciiSubscriptionsStats(team.ChallengeSubscriptions)
					achievements := fmt.Sprintf("%d", len(team.Achievements))
					slug := team.Slug
					cash := fmt.Sprintf("$%d", team.Cash)
					if team.Cash == 0 {
						cash = "🚫"
					}
					table.Append([]string{slug, createdAgo, updatedAgo, season, cash, members, organization, challenges, achievements, status, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminCouponsCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin coupons", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "coupons",
		ShortUsage: "pathwar [global flags] admin [admin flags] coupons [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListCoupons(ctx, &pwapi.AdminListCoupons_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// coupons table
			{
				fmt.Println("COUPONS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"HASH", "VALUE", "SEASON", "CREATED", "UPDATED", "VALIDATIONS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, coupon := range ret.Coupons {
					// fmt.Println(godev.PrettyJSONPB(coupon))
					hash := coupon.Hash
					id := fmt.Sprintf("%d", coupon.ID)
					value := fmt.Sprintf("%d", coupon.Value)
					createdAgo := humanize.Time(*coupon.CreatedAt)
					updatedAgo := humanize.Time(*coupon.UpdatedAt)
					season := coupon.Season.ASCIIID()
					validated := int64(len(coupon.Validations))
					validationStatus := "🟢"
					switch {
					case validated > 0 && validated < coupon.MaxValidationCount:
						validationStatus = "🔶"
					case validated == coupon.MaxValidationCount:
						validationStatus = "🔴"
					case validated > coupon.MaxValidationCount:
						validationStatus = "🙀"
					}
					validations := fmt.Sprintf("%d / %d %s", validated, coupon.MaxValidationCount, validationStatus)
					table.Append([]string{hash, value, season, createdAgo, updatedAgo, validations, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminChallengeSubscriptionsCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin challengeSubscriptions", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "subscriptions",
		ShortUsage: "pathwar [global flags] admin [admin flags] subscriptions [flags]",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListChallengeSubscriptions(ctx, &pwapi.AdminListChallengeSubscriptions_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			// challengeSubscriptions table
			{
				fmt.Println("CHALLENGE SUBSCRIPTIONS")
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"CHALLENGE", "TEAM", "SEASON", "STATUS", "CREATED", "UPDATED", "BUYER", "CLOSER", "VALIDATIONS", "ID"})
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				table.SetBorder(false)

				for _, subscription := range ret.Subscriptions {
					// fmt.Println(godev.PrettyJSONPB(subscription))
					id := fmt.Sprintf("%d", subscription.ID)
					createdAgo := humanize.Time(*subscription.CreatedAt)
					updatedAgo := humanize.Time(*subscription.UpdatedAt)
					team := subscription.Team.Organization.ASCIIID()
					buyer := subscription.Buyer.ASCIIID()
					season := subscription.SeasonChallenge.Season.ASCIIID()
					challenge := subscription.SeasonChallenge.Flavor.Challenge.ASCIIID()
					status := asciiStatus(subscription.Status.String())
					closer := subscription.Closer.ASCIIID()
					validations := fmt.Sprintf("%d", len(subscription.Validations))
					table.Append([]string{challenge, team, season, status, createdAgo, updatedAgo, buyer, closer, validations, id})
				}
				table.Render()
				fmt.Println("")
			}

			return nil
		},
	}
}

func adminListAllCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "list-all",
		ShortUsage: "pathwar [global flags] admin [admin flags] list-all",
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminListAll(ctx, &pwapi.AdminListAll_Input{})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}
			fmt.Println("LIST ALL")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "SLUG", "type"})
			table.SetAlignment(tablewriter.ALIGN_CENTER)
			table.SetBorder(false)

			for _, challenge := range ret.Challenges {
				table.Append([]string{fmt.Sprint(challenge.ID), challenge.Slug, "challenge"})
			}
			for _, challengeFlavor := range ret.ChallengeFlavors {
				table.Append([]string{fmt.Sprint(challengeFlavor.ID), challengeFlavor.Slug, "challenge-flavor"})
			}
			for _, seasonChallenge := range ret.SeasonChallenges {
				table.Append([]string{fmt.Sprint(seasonChallenge.ID), seasonChallenge.Slug, "season-challenge"})
			}
			for _, challengeInstance := range ret.ChallengeInstances {
				table.Append([]string{fmt.Sprint(challengeInstance.ID), challengeInstance.Slug, "challenge-instance"})
			}
			for _, agent := range ret.Agents {
				table.Append([]string{fmt.Sprint(agent.ID), agent.Slug, "agent"})
			}
			for _, organizationMember := range ret.OrganizationMembers {
				table.Append([]string{fmt.Sprint(organizationMember.ID), organizationMember.Slug, "organization-member"})
			}
			for _, teamMember := range ret.TeamMembers {
				table.Append([]string{fmt.Sprint(teamMember.ID), teamMember.Slug, "team-member"})
			}
			for _, teamInvite := range ret.TeamInvites {
				table.Append([]string{fmt.Sprint(teamInvite.ID), teamInvite.Slug, "team-invite"})
			}
			for _, user := range ret.Users {
				table.Append([]string{fmt.Sprint(user.ID), user.Slug, "user"})
			}
			for _, organization := range ret.Organizations {
				table.Append([]string{fmt.Sprint(organization.ID), organization.Slug, "organization"})
			}
			for _, season := range ret.Seasons {
				table.Append([]string{fmt.Sprint(season.ID), season.Slug, "season"})
			}
			for _, team := range ret.Teams {
				table.Append([]string{fmt.Sprint(team.ID), team.Slug, "team"})
			}
			for _, whoswhoAttempt := range ret.WhoswhoAttempts {
				table.Append([]string{fmt.Sprint(whoswhoAttempt.ID), whoswhoAttempt.Slug, "whoswho-attempt"})
			}
			for _, challengeValidation := range ret.ChallengeValidations {
				table.Append([]string{fmt.Sprint(challengeValidation.ID), challengeValidation.Slug, "challenge-validation"})
			}
			for _, challengeSubscription := range ret.ChallengeSubscriptions {
				table.Append([]string{fmt.Sprint(challengeSubscription.ID), challengeSubscription.Slug, "challenge-subscription"})
			}
			for _, inventoryItem := range ret.InventoryItems {
				table.Append([]string{fmt.Sprint(inventoryItem.ID), inventoryItem.Slug, "inventory-item"})
			}
			/*for _, notification := range ret.Notifications {
				table.Append([]string{fmt.Sprint(notification.ID), notification.Slug, "notification"})
			}
			for _, coupon := range ret.Coupons {
				table.Append([]string{fmt.Sprint(coupon.ID), coupon.Slug, "coupon"})
			}
			for _, couponValidation := range ret.CouponValidations {
				table.Append([]string{fmt.Sprint(couponValidation.ID), couponValidation.Slug, "coupon-validation"})
			}
			for _, achievement := range ret.Achievements {
				table.Append([]string{fmt.Sprint(achievement.ID), achievement.Slug, "achievement"})
			}
			for _, activity := range ret.Activities {
				table.Append([]string{fmt.Sprint(activity.ID), activity.Slug, "activity"})
			}*/

			table.Render()
			fmt.Println("")
			return nil
		},
	}
}

func adminSearchCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "search",
		ShortUsage: "pathwar [global flags] admin [admin flags] search QUERY",
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			if len(args) != 1 {
				return flag.ErrHelp
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminSearch(ctx, &pwapi.AdminSearch_Input{Search: args[0]})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			fmt.Println(godev.PrettyJSONPB(&ret))
			return nil
		},
	}
}

func adminAddCouponCommand() *ffcli.Command {
	input := pwapi.AdminAddCoupon_Input{}
	input.ApplyDefaults()

	flags := flag.NewFlagSet("admin add-coupon", flag.ExitOnError)
	flags.StringVar(&input.Hash, "hash", input.Hash, "Hash to guess (must be unique, if 'RANDOM', will be randomized)")
	flags.StringVar(&input.SeasonID, "season", input.SeasonID, "Season ID or Slug to associate the coupon with")
	flags.Int64Var(&input.Value, "value", input.Value, "Coupon value")
	flags.Int64Var(&input.MaxValidationCount, "max-validations", input.MaxValidationCount, "Maximum times a coupon can be validated")
	return &ffcli.Command{
		Name:       "add-coupon",
		ShortUsage: "pathwar admin add-coupon",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminAddCoupon(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
			} else {
				fmt.Println(ret.Coupon.Hash)
			}
			return nil
		},
	}
}

func adminRedumpCommand() *ffcli.Command {
	flags := flag.NewFlagSet("admin redump", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "redump",
		ShortUsage: "pathwar [global flags] admin [admin flags] redump [flags] ID...",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminRedump(ctx, &pwapi.AdminRedump_Input{
				Identifiers: args,
			})
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println("OK")

			return nil
		},
	}
}

func adminChallengeAddCommand() *ffcli.Command {
	input := pwapi.AdminChallengeAdd_Input{Challenge: &pwdb.Challenge{}}
	input.ApplyDefaults()
	flags := flag.NewFlagSet("admin challenge add", flag.ExitOnError)
	flags.StringVar(&input.Challenge.Slug, "slug", input.Challenge.Slug, "Unique slug")
	flags.StringVar(&input.Challenge.Name, "name", input.Challenge.Name, "Challenge name")
	flags.StringVar(&input.Challenge.Description, "description", input.Challenge.Description, "Challenge description")
	flags.StringVar(&input.Challenge.Author, "author", input.Challenge.Author, "Challenge author")
	flags.StringVar(&input.Challenge.Locale, "locale", input.Challenge.Locale, "Challenge Locale")
	flags.BoolVar(&input.Challenge.IsDraft, "is-draft", input.Challenge.IsDraft, "Is challenge production ready ?")
	flags.StringVar(&input.Challenge.PreviewURL, "preview-url", input.Challenge.PreviewURL, "Challenge preview URL")
	flags.StringVar(&input.Challenge.Homepage, "homepage", input.Challenge.Homepage, "Challenge homepage URL")

	return &ffcli.Command{
		Name:       "challenge-add",
		ShortUsage: "pathwar [global flags] admin [admin flags] challenge-add [flags] [args...]",
		ShortHelp:  "add a challenge",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if input.Challenge.Name == "" {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminAddChallenge(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.Challenge.ID)
			return nil
		},
	}
}

func adminChallengeRedumpCommand() *ffcli.Command {
	input := pwapi.AdminChallengeRedump_Input{}
	flags := flag.NewFlagSet("admin challenge redump", flag.ExitOnError)
	flags.StringVar(&input.ChallengeID, "id", input.ChallengeID, "Challenge ID or slug")

	return &ffcli.Command{
		Name:       "challenge-redump",
		ShortUsage: "pathwar [global flags] admin [admin flags] challenge-redump [flags] [args...]",
		ShortHelp:  "redump a challenge",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if input.ChallengeID == "" {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminRedumpChallenge(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.ChallengeInstances)
			return nil
		},
	}
}

func adminChallengeFlavorAddCommand() *ffcli.Command {
	input := pwapi.AdminChallengeFlavorAdd_Input{}
	input.ApplyDefaults()
	flags := flag.NewFlagSet("admin challenge flavor add", flag.ExitOnError)
	flags.StringVar(&input.ChallengeID, "challenge", input.ChallengeID, "Challenge ID or slug")
	flags.StringVar(&input.ChallengeFlavor.Slug, "slug", input.ChallengeFlavor.Slug, "Slug")
	flags.StringVar(&input.ChallengeFlavor.Version, "version", input.ChallengeFlavor.Version, "Challenge flavor version")
	flags.StringVar(&input.ChallengeFlavor.ComposeBundle, "compose-bundle", input.ChallengeFlavor.ComposeBundle, "Challenge flavor compose bundle")
	flags.StringVar(&input.ChallengeFlavor.SourceURL, "source-url", input.ChallengeFlavor.SourceURL, "Source URL")
	flags.Int64Var(&input.ChallengeFlavor.PurchasePrice, "purchase-price", input.ChallengeFlavor.PurchasePrice, "Purchase Price")
	flags.Int64Var(&input.ChallengeFlavor.ValidationReward, "validation-reward", input.ChallengeFlavor.ValidationReward, "Validation reward")
	flags.StringVar(&input.ChallengeFlavor.Body, "body", input.ChallengeFlavor.Body, "Body")
	flags.StringVar(&input.ChallengeFlavor.Category, "category", input.ChallengeFlavor.Category, "Category")
	flags.StringVar(&input.ChallengeFlavor.TagList, "tags", input.ChallengeFlavor.TagList, "Comma-separated tags")
	flags.StringVar(&input.ChallengeFlavor.RedumpPolicyConfig, "redump-policy", input.ChallengeFlavor.RedumpPolicyConfig, "JSON config for redump-policy")
	flags.Int64Var(&input.ChallengeFlavor.Passphrases, "passphrases", input.ChallengeFlavor.Passphrases, "Amount of passphrases")

	return &ffcli.Command{
		Name:       "challenge-flavor-add",
		ShortUsage: "pathwar [global flags] admin [admin flags] challenge-flavor-add [flags] [args...]",
		ShortHelp:  "add a challenge flavor",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			input.ApplyDefaults()
			if input.ChallengeID == "" {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			if compose := input.ChallengeFlavor.ComposeBundle; compose != "" {
				composePath, err := filepath.Abs(compose)
				if err != nil {
					return errcode.TODO.Wrap(err)
				}
				if _, err := os.Stat(composePath); err == nil {
					f, err := os.Open(composePath)
					if err != nil {
						return errcode.TODO.Wrap(err)
					}
					defer f.Close()
					b, err := ioutil.ReadAll(f)
					if err != nil {
						return errcode.TODO.Wrap(err)
					}
					input.ChallengeFlavor.ComposeBundle = string(b)
				} else {
					return errcode.TODO.Wrap(err)
				}
			}

			ret, err := apiClient.AdminAddChallengeFlavor(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.ChallengeFlavor.ID)
			return nil
		},
	}
}

func adminSeasonAddCommand() *ffcli.Command {
	input := pwapi.AdminSeasonAdd_Input{Season: &pwdb.Season{}}
	visibility := "private"
	status := "stopped"
	subscription := "close"
	flags := flag.NewFlagSet("admin season add", flag.ExitOnError)
	flags.StringVar(&input.Season.Name, "name", input.Season.Name, "Name")
	flags.BoolVar(&input.Season.IsGlobal, "global", input.Season.IsGlobal, "Global season ?")
	flags.BoolVar(&input.Season.IsTesting, "testing", input.Season.IsGlobal, "Testing season ?")
	flags.StringVar(&visibility, "visibility", visibility, "Season visibility (private, unlisted or public)")
	flags.StringVar(&status, "status", status, "Season status (started or stopped)")
	flags.StringVar(&subscription, "subscription", subscription, "Susbscription status (open or close)")

	return &ffcli.Command{
		Name:       "season-add",
		ShortUsage: "pathwar [global flags] admin [admin flags] season-add [flags] [args...]",
		ShortHelp:  "add a Season",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if input.Season.Name == "" {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			switch visibility {
			case "unlisted":
				input.Season.Visibility = pwdb.Season_Unlisted
			case "public":
				input.Season.Visibility = pwdb.Season_Public
			case "private":
				input.Season.Visibility = pwdb.Season_Private
			default:
				return fmt.Errorf("unsupported visibility: %q", visibility)
			}
			switch status {
			case "started":
				input.Season.Status = pwdb.Season_Started
			case "stopped":
				input.Season.Status = pwdb.Season_Stopped
			default:
				return fmt.Errorf("unsupported status: %q", status)
			}
			switch subscription {
			case "open":
				input.Season.Subscription = pwdb.Season_Open
			case "close":
				input.Season.Subscription = pwdb.Season_Close
			default:
				return fmt.Errorf("unsupported subscription: %q", subscription)
			}

			ret, err := apiClient.AdminAddSeason(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.Season.Slug)
			return nil
		},
	}
}

func adminSeasonChallengeAddCommand() *ffcli.Command {
	input := pwapi.AdminSeasonChallengeAdd_Input{}
	input.ApplyDefaults()
	flags := flag.NewFlagSet("admin season challenge add", flag.ExitOnError)
	flags.StringVar(&input.FlavorID, "flavor", input.FlavorID, "Flavor ID or Slug")
	flags.StringVar(&input.SeasonID, "season", input.SeasonID, "Season ID or Slug")
	flags.StringVar(&input.SeasonChallenge.Slug, "slug", input.SeasonChallenge.Slug, "Slug")

	return &ffcli.Command{
		Name:       "season-challenge-add",
		ShortUsage: "pathwar [global flags] admin [admin flags] season-challenge-add [flags] [args...]",
		ShortHelp:  "add a SeasonChallenge",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			input.ApplyDefaults()
			if input.FlavorID == "" || input.SeasonID == "" {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminAddSeasonChallenge(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.SeasonChallenge.ID)
			return nil
		},
	}
}

func adminTestingSeasonUserAddCommand() *ffcli.Command {
	input := pwapi.AdminTestingSeasonUserAdd_Input{}
	flags := flag.NewFlagSet("admin season testing user add", flag.ExitOnError)
	flags.StringVar(&input.UserID, "user", input.UserID, "User ID or Slug")

	return &ffcli.Command{
		Name:       "testing-season-user-add",
		ShortUsage: "pathwar [global flags] admin [admin flags] testing-season-user-add [flags] [args...]",
		ShortHelp:  "add a user to testing season",
		FlagSet:    flags,
		Exec: func(ctx context.Context, args []string) error {
			if input.UserID == "" {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			ret, err := apiClient.AdminTestingSeasonUserAdd(ctx, &input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}
			if globalDebug {
				fmt.Fprintln(os.Stderr, godev.PrettyJSONPB(&ret))
			}

			if adminJSONFormat {
				fmt.Println(godev.PrettyJSONPB(&ret))
				return nil
			}

			fmt.Println(ret.TeamMember.ID)
			return nil
		},
	}
}

func asciiInstancesStats(instances []*pwdb.ChallengeInstance) string {
	if len(instances) == 0 {
		return "🚫"
	}

	instanceGreen := 0
	instanceRed := 0
	for _, instance := range instances {
		if instance.Status == pwdb.ChallengeInstance_Available {
			instanceGreen++
		} else {
			instanceRed++
		}
	}
	instanceParts := []string{}
	if instanceGreen > 0 {
		instanceParts = append(instanceParts, fmt.Sprintf("%dx🟢", instanceGreen))
	}
	if instanceRed > 0 {
		instanceParts = append(instanceParts, fmt.Sprintf("%dx🔴", instanceRed))
	}
	stats := strings.Join(instanceParts, " + ")
	return stats
}

func asciiSubscriptionsStats(subscriptions []*pwdb.ChallengeSubscription) string {
	if len(subscriptions) == 0 {
		return "🚫"
	}

	subscriptionGreen := 0
	subscriptionRed := 0
	for _, subscription := range subscriptions {
		if subscription.Status == pwdb.ChallengeSubscription_Active {
			subscriptionGreen++
		} else {
			subscriptionRed++
		}
	}
	subscriptionParts := []string{}
	if subscriptionGreen > 0 {
		subscriptionParts = append(subscriptionParts, fmt.Sprintf("%dx🟢", subscriptionGreen))
	}
	if subscriptionRed > 0 {
		subscriptionParts = append(subscriptionParts, fmt.Sprintf("%dx🔴", subscriptionRed))
	}
	stats := strings.Join(subscriptionParts, " + ")
	return stats
}

func asciiBool(input bool) string {
	if !input {
		return "❌"
	}
	return "✅"
}

func asciiStatus(status string) string {
	switch status {
	case "Active", "Available":
		status += " 🟢"
	default:
		status += " 🔴"
	}
	return status
}
