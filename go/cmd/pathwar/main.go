package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"

	bearer "github.com/Bearer/bearer-go"
	"github.com/bwmarrin/snowflake"
	"github.com/docker/docker/client"
	sentry "github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql" // required by gorm
	"github.com/jinzhu/gorm"
	"github.com/oklog/run"
	opentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"moul.io/banner"
	"moul.io/godev"
	"moul.io/motd"
	"moul.io/srand"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwagent"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
	"pathwar.land/pathwar/v2/go/pkg/pwversion"
)

const (
	defaultDBURN       = "root:uns3cur3@tcp(127.0.0.1:3306)/pathwar?charset=utf8mb4&parseTime=true"
	defaultHTTPApiAddr = "https://api-dev.pathwar.land"
)

var (
	logger *zap.Logger
	tracer opentracing.Tracer

	flagOutput = os.Stderr

	// flag vars
	adminChallengeAddInput         = pwapi.AdminChallengeAdd_Input{Challenge: &pwdb.Challenge{}}
	adminChallengeFlavorAddInput   = pwapi.AdminChallengeFlavorAdd_Input{ChallengeFlavor: &pwdb.ChallengeFlavor{}}
	adminChallengeInstanceAddInput = pwapi.AdminChallengeInstanceAdd_Input{ChallengeInstance: &pwdb.ChallengeInstance{}}
	agentOpts                      = pwagent.NewOpts()
	serverOpts                     = pwapi.NewServerOpts()
	ssoOpts                        = pwsso.NewOpts()
	composeCleanOpts               = pwcompose.NewCleanOpts()
	composePrepareOpts             = pwcompose.NewPrepareOpts()
	composeUpOpts                  = pwcompose.NewUpOpts()

	DBURN           string
	DBMaxOpenTries  int
	bearerSecretKey string
	composePSDepth  int
	globalDebug     bool
	globalSentryDSN string
	httpAPIAddr     string
	zipkinEndpoint  string
)

// nolint:gocyclo
func main() {
	log.SetFlags(0)

	defer func() {
		if logger != nil {
			_ = logger.Sync()
		}
	}()

	// setup flags
	var (
		globalFlags                    = flag.NewFlagSet("pathwar", flag.ExitOnError)
		agentFlags                     = flag.NewFlagSet("agent", flag.ExitOnError)
		clientFlags                    = flag.NewFlagSet("client", flag.ExitOnError)
		apiFlags                       = flag.NewFlagSet("api", flag.ExitOnError)
		composeDownFlags               = flag.NewFlagSet("compose down", flag.ExitOnError)
		composeFlags                   = flag.NewFlagSet("compose", flag.ExitOnError)
		composePSFlags                 = flag.NewFlagSet("compose ps", flag.ExitOnError)
		composePrepareFlags            = flag.NewFlagSet("compose prepare", flag.ExitOnError)
		composeUpFlags                 = flag.NewFlagSet("compose up", flag.ExitOnError)
		miscFlags                      = flag.NewFlagSet("misc", flag.ExitOnError)
		serverFlags                    = flag.NewFlagSet("server", flag.ExitOnError)
		ssoFlags                       = flag.NewFlagSet("sso", flag.ExitOnError)
		adminFlags                     = flag.NewFlagSet("admin", flag.ExitOnError)
		adminPSFlags                   = flag.NewFlagSet("admin ps", flag.ExitOnError)
		adminRedumpFlags               = flag.NewFlagSet("admin redump", flag.ExitOnError)
		adminChallengeAddFlags         = flag.NewFlagSet("admin challenge add", flag.ExitOnError)
		adminChallengeFlavorAddFlags   = flag.NewFlagSet("admin challenge flavor add", flag.ExitOnError)
		adminChallengeInstanceAddFlags = flag.NewFlagSet("admin challenge instance add", flag.ExitOnError)
	)
	globalFlags.SetOutput(flagOutput) // used in main_test.go
	globalFlags.BoolVar(&globalDebug, "debug", false, "debug mode")
	globalFlags.StringVar(&zipkinEndpoint, "zipkin-endpoint", "", "optional opentracing server")
	globalFlags.StringVar(&bearerSecretKey, "bearer-secretkey", "", "bearer.sh secret key")
	globalFlags.StringVar(&globalSentryDSN, "sentry-dsn", "", "Sentry DSN")

	agentFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	agentFlags.StringVar(&ssoOpts.ClientID, "sso-clientid", ssoOpts.ClientID, "SSO ClientID")
	agentFlags.StringVar(&ssoOpts.ClientSecret, "sso-clientsecret", ssoOpts.ClientSecret, "SSO ClientSecret")
	agentFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")
	agentFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")

	agentFlags.BoolVar(&agentOpts.Cleanup, "clean", agentOpts.Cleanup, "remove all pathwar instances before executing")
	agentFlags.BoolVar(&agentOpts.RunOnce, "once", agentOpts.RunOnce, "run once and don't start daemon loop")
	agentFlags.BoolVar(&agentOpts.NoRun, "no-run", agentOpts.NoRun, "stop after agent initialization (register and cleanup)")
	agentFlags.DurationVar(&agentOpts.LoopDelay, "delay", agentOpts.LoopDelay, "delay between each loop iteration")
	agentFlags.BoolVar(&agentOpts.DefaultAgent, "default-agent", agentOpts.DefaultAgent, "agent creates an instance for each available flavor on registration, else will only create an instance of debug-challenge")
	agentFlags.StringVar(&agentOpts.Name, "agent-name", agentOpts.Name, "Agent Name")
	agentFlags.StringVar(&agentOpts.DomainSuffix, "domain-suffix", agentOpts.DomainSuffix, "Domain suffix to append")
	agentFlags.StringVar(&agentOpts.NginxDockerImage, "docker-image", agentOpts.NginxDockerImage, "docker image used to generate nginx proxy container")
	agentFlags.StringVar(&agentOpts.HostIP, "host", agentOpts.HostIP, "Nginx HTTP listening addr")
	agentFlags.StringVar(&agentOpts.HostPort, "port", agentOpts.HostPort, "Nginx HTTP listening port")
	agentFlags.StringVar(&agentOpts.ModeratorPassword, "moderator-password", agentOpts.ModeratorPassword, "Challenge moderator password")
	agentFlags.StringVar(&agentOpts.AuthSalt, "salt", agentOpts.AuthSalt, "salt used to generate secure hashes (random if empty)")

	adminFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	adminFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")

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

	apiFlags.BoolVar(&ssoOpts.AllowUnsafe, "sso-unsafe", ssoOpts.AllowUnsafe, "Allow unsafe SSO")
	apiFlags.StringVar(&DBURN, "urn", defaultDBURN, "MySQL URN")
	apiFlags.IntVar(&DBMaxOpenTries, "db-max-open-tries", 0, "max DB open tries, unlimited if 0")
	apiFlags.StringVar(&ssoOpts.ClientID, "sso-clientid", ssoOpts.ClientID, "SSO ClientID")
	apiFlags.StringVar(&ssoOpts.Pubkey, "sso-pubkey", ssoOpts.Pubkey, "SSO Public Key")
	apiFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")

	clientFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	clientFlags.StringVar(&ssoOpts.ClientID, "sso-clientid", ssoOpts.ClientID, "SSO ClientID")
	clientFlags.StringVar(&ssoOpts.ClientSecret, "sso-clientsecret", ssoOpts.ClientSecret, "SSO ClientSecret")
	clientFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")
	clientFlags.StringVar(&ssoOpts.TokenFile, "sso-token-file", ssoOpts.TokenFile, "Token file")

	composeDownFlags.BoolVar(&composeCleanOpts.RemoveVolumes, "rm-volumes", composeCleanOpts.RemoveVolumes, "keep volumes")
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveImages, "rm-images", composeCleanOpts.RemoveImages, "remove images as well")
	composeDownFlags.BoolVar(&composeCleanOpts.RemoveNginx, "rm-nginx", composeCleanOpts.RemoveNginx, "down nginx container and proxy network as well")

	composePSFlags.IntVar(&composePSDepth, "depth", 0, "depth to display")

	composePrepareFlags.BoolVar(&composePrepareOpts.NoPush, "no-push", composePrepareOpts.NoPush, "don't push images")
	composePrepareFlags.StringVar(&composePrepareOpts.Prefix, "prefix", composePrepareOpts.Prefix, "docker image prefix")
	composePrepareFlags.StringVar(&composePrepareOpts.Version, "version", composePrepareOpts.Version, "challenge version")

	composeUpFlags.StringVar(&composeUpOpts.InstanceKey, "instance-key", composeUpOpts.InstanceKey, "instance key used to generate instance ID")
	composeUpFlags.BoolVar(&composeUpOpts.ForceRecreate, "force-recreate", composeUpOpts.ForceRecreate, "down previously created instances of challenge")

	serverFlags.BoolVar(&serverOpts.WithPprof, "with-pprof", serverOpts.WithPprof, "enable pprof endpoints")
	serverFlags.DurationVar(&serverOpts.RequestTimeout, "request-timeout", serverOpts.RequestTimeout, "request timeout")
	serverFlags.DurationVar(&serverOpts.ShutdownTimeout, "shutdown-timeout", serverOpts.ShutdownTimeout, "shutdown timeout")
	serverFlags.StringVar(&serverOpts.CORSAllowedOrigins, "cors-allowed-origins", serverOpts.CORSAllowedOrigins, "allowed CORS origins")
	serverFlags.StringVar(&serverOpts.Bind, "bind", serverOpts.Bind, "server address")

	ssoFlags.BoolVar(&ssoOpts.AllowUnsafe, "unsafe", ssoOpts.AllowUnsafe, "Allow unsafe SSO")
	ssoFlags.StringVar(&ssoOpts.ClientID, "clientid", ssoOpts.ClientID, "SSO ClientID")
	ssoFlags.StringVar(&ssoOpts.Pubkey, "pubkey", ssoOpts.Pubkey, "SSO Public Key")
	ssoFlags.StringVar(&ssoOpts.Realm, "realm", ssoOpts.Realm, "SSO Realm")

	version := &ffcli.Command{
		Name:      "version",
		Usage:     "pathwar [global flags] version",
		ShortHelp: "show version",
		Exec: func(args []string) error {
			fmt.Printf(
				"version=%q\ncommit=%q\nbuilt-at=%q\nbuilt-by=%q\n",
				pwversion.Version, pwversion.Commit, pwversion.Date, pwversion.BuiltBy,
			)
			return nil
		},
	}

	api := &ffcli.Command{
		Name:      "api",
		Usage:     "pathwar [global flags] api [api flags] <subcommand> [flags] [args...]",
		ShortHelp: "manage the Pathwar API",
		FlagSet:   apiFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Subcommands: []*ffcli.Command{{
			Name:      "server",
			Usage:     "pathwar [global flags] server [server flags] <subcommand> [flags] [args...]",
			ShortHelp: "start a server (HTTP + gRPC)",
			FlagSet:   serverFlags,
			Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(args []string) error {
				fmt.Println(motd.Default())
				fmt.Println(banner.Inline("api server"))

				err := globalPreRun()
				if err != nil {
					return err
				}
				cleanup, err := initSentryFromEnv("starting API")
				if err != nil {
					return err
				}
				defer cleanup()

				// init svc
				svc, _, closer, err := svcFromFlags(logger)
				if err != nil {
					return errcode.ErrStartService.Wrap(err)
				}
				defer closer()

				var (
					ctx    = context.Background()
					g      run.Group
					server *pwapi.Server
				)
				g.Add(run.SignalHandler(ctx, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill))
				{ // server
					serverOpts.Tracer = tracer
					serverOpts.Logger = logger.Named("server")
					var err error
					server, err = pwapi.NewServer(ctx, svc, serverOpts)
					if err != nil {
						return errcode.ErrInitServer.Wrap(err)
					}
					g.Add(
						server.Run,
						func(error) { server.Close() },
					)
				}

				logger.Info("server started",
					zap.String("bind", server.ListenerAddr()),
				)

				if err := g.Run(); err != nil {
					return errcode.ErrGroupTerminated.Wrap(err)
				}
				return nil
			},
		}, {
			Name:  "sqlinfo",
			Usage: "pathwar [global flags] api [api flags] sqlinfo",
			Exec: func([]string) error {
				err := globalPreRun()
				if err != nil {
					return err
				}

				// init svc
				_, db, closer, err := svcFromFlags(logger)
				if err != nil {
					return errcode.ErrStartService.Wrap(err)
				}
				defer closer()

				info, err := pwdb.GetInfo(db, logger)
				if err != nil {
					return errcode.ErrGetDBInfo.Wrap(err)
				}

				out, _ := json.MarshalIndent(info, "", "  ")
				fmt.Println(string(out))
				return nil
			},
		}, {
			Name:  "sqldump",
			Usage: "pathwar [global flags] api [api flags] sqldump",
			Exec: func([]string) error {
				err := globalPreRun()
				if err != nil {
					return err
				}

				// init svc
				_, db, closer, err := svcFromFlags(logger)
				if err != nil {
					return errcode.ErrStartService.Wrap(err)
				}
				defer closer()

				dump, err := pwdb.GetDump(db)
				if err != nil {
					return errcode.ErrDumpDatabase.Wrap(err)
				}

				out, _ := json.MarshalIndent(dump, "", "  ")
				fmt.Println(string(out))
				return nil
			},
		}},
		Exec: func([]string) error { return flag.ErrHelp },
	}

	misc := &ffcli.Command{
		Name:      "misc",
		Usage:     "pathwar [global flags] misc [misc flags] <subcommand> [flags] [args...]",
		ShortHelp: "misc contains advanced commands",
		Subcommands: []*ffcli.Command{{
			Name:  "pwinit-binary",
			Usage: "pathwar [global flags] misc [misc flags] pwinit-binary",
			Exec: func([]string) error {
				binary, err := pwinit.Binary()
				if err != nil {
					return err
				}
				os.Stdout.Write(binary)
				return nil
			},
		}},
		FlagSet: miscFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:    func([]string) error { return flag.ErrHelp },
	}

	sso := &ffcli.Command{
		Name:      "sso",
		Usage:     "pathwar [global flags] sso [sso flags] <subcommand> [flags] [args...]",
		ShortHelp: "manage SSO tokens",
		Subcommands: []*ffcli.Command{{
			Name:  "token",
			Usage: "pathwar [global flags] sso [sso flags] token TOKEN",
			Exec: func(args []string) error {
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
			Name:  "logout",
			Usage: "pathwar [global flags] sso [sso flags] logout TOKEN",
			Exec: func(args []string) error {
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
			Name:  "whoami",
			Usage: "pathwar [global flags] sso [sso flags] whoami TOKEN",
			Exec: func(args []string) error {
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
		}},
		FlagSet: ssoFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:    func([]string) error { return flag.ErrHelp },
	}

	compose := &ffcli.Command{
		Name:  "compose",
		Usage: "pathwar [global flags] compose [compose flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{{
			Name:    "up",
			Usage:   "pathwar [global flags] compose [compose flags] up [flags] PATH",
			FlagSet: composeUpFlags,
			Options: []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(args []string) error {
				err := globalPreRun()
				if err != nil {
					return err
				}

				if len(args) < 1 {
					return flag.ErrHelp
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
				if len(args) < 1 {
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

	admin := &ffcli.Command{
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

				fmt.Println(godev.PrettyJSONPB(&ret))
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

	clientCmd := &ffcli.Command{
		Name:      "client",
		Usage:     "pathwar [global flags] client [client flags] <method> <path> [INPUT (json)]",
		ShortHelp: "make API calls",
		LongHelp: `EXAMPLES
  pathwar client GET /user/session
  season=$(pathwar client GET /user/session | jq -r '.seasons[0].season.id')
  pathwar client GET "/season-challenges?season_id=$season"`,
		FlagSet: clientFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec: func(args []string) error {
			if len(args) < 2 || len(args) > 3 {
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

			method := args[0]
			path := args[1]
			var input []byte
			if len(args) > 2 {
				input = []byte(args[2])
			}
			output, err := apiClient.Raw(ctx, method, path, input)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			var data interface{}
			err = json.Unmarshal(output, &data)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			fmt.Println(godev.PrettyJSON(data))
			return nil
		},
	}

	agent := &ffcli.Command{
		Name:      "agent",
		Usage:     "pathwar [global flags] agent [agent flags] <subcommand> [flags] [args...]",
		ShortHelp: "manage an agent node (multiple challenges)",
		FlagSet:   agentFlags,
		Options:   []ff.Option{ff.WithEnvVarNoPrefix()},
		Subcommands: []*ffcli.Command{
			{
				Name:      "pwinit.bin",
				ShortHelp: "dump pwinit binary to stdout",
				Exec: func(args []string) error {
					b, err := pwinit.Binary()
					if err != nil {
						return err
					}
					_, err = os.Stdout.Write(b)
					return err
				},
			},
		},
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			fmt.Println(motd.Default())
			fmt.Println(banner.Inline("agent"))

			cleanup, err := initSentryFromEnv("starting agent")
			if err != nil {
				return err
			}
			defer cleanup()

			ctx := context.Background()
			dockerCli, err := client.NewEnvClient()
			if err != nil {
				return errcode.ErrInitDockerClient.Wrap(err)
			}

			apiClient, err := httpClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			agentOpts.Logger = logger
			return pwagent.Run(ctx, dockerCli, apiClient, agentOpts)
		},
	}

	root := &ffcli.Command{
		Usage:       "pathwar [global flags] <subcommand> [flags] [args...]",
		FlagSet:     globalFlags,
		LongHelp:    "More info here: https://github.com/pathwar/pathwar/wiki/CLI",
		Subcommands: []*ffcli.Command{clientCmd, api, compose, agent, sso, misc, version, admin},
		Options:     []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	if err := root.Run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatalf("fatal: %+v", err)
	}
}

func svcFromFlags(logger *zap.Logger) (pwapi.Service, *gorm.DB, func(), error) {
	// init database
	dbConnectTries := 0
dbConnectLoop:
	db, err := gorm.Open("mysql", DBURN)
	if err != nil {
		dbConnectTries++
		fmt.Println(DBMaxOpenTries, dbConnectTries)
		if DBMaxOpenTries == 0 || dbConnectTries < DBMaxOpenTries {
			logger.Warn("db open", zap.Int("tries", dbConnectTries), zap.Error(err))
			time.Sleep(5 * time.Second)
			goto dbConnectLoop
		}
		return nil, nil, nil, errcode.ErrInitDB.Wrap(err)
	}
	sfn, err := snowflake.NewNode(1)
	if err != nil {
		return nil, nil, nil, errcode.ErrInitSnowflake.Wrap(err)
	}
	dbOpts := pwdb.Opts{
		Logger: logger.Named("gorm"),
	}
	db, err = pwdb.Configure(db, sfn, dbOpts)
	if err != nil {
		return nil, nil, nil, errcode.ErrConfigureDB.Wrap(err)
	}

	// init SSO
	sso, err := ssoFromFlags()
	if err != nil {
		return nil, nil, nil, errcode.ErrInitSSOClient.Wrap(err)
	}

	// init svc
	svcOpts := pwapi.ServiceOpts{
		Logger: logger.Named("svc"),
	}

	svc, err := pwapi.NewService(db, sso, svcOpts)
	if err != nil {
		return nil, nil, nil, errcode.ErrInitService.Wrap(err)
	}

	closeFunc := func() {
		if err := svc.Close(); err != nil {
			logger.Warn("close svc", zap.Error(err))
		}
		if err := db.Close(); err != nil {
			logger.Warn("closed database", zap.Error(err))
		}
	}

	logger.Debug("svc initd", zap.Any("db", db), zap.Any("opts", svcOpts))
	return svc, db, closeFunc, nil
}

func ssoFromFlags() (pwsso.Client, error) {
	ssoOpts.Logger = logger.Named("sso")
	ssoOpts.ApplyDefaults()
	sso, err := pwsso.New(ssoOpts.Pubkey, ssoOpts.Realm, ssoOpts)
	if err != nil {
		return nil, errcode.ErrInitSSOClient.Wrap(err)
	}
	return sso, nil
}

func globalPreRun() error {
	rand.Seed(srand.Secure())
	if bearerSecretKey != "" {
		bearer.ReplaceGlobals(bearer.Init(bearerSecretKey))
	}
	if globalDebug {
		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(zap.DebugLevel)
		config.DisableStacktrace = true
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		var err error
		logger, err = config.Build()
		if err != nil {
			return errcode.ErrInitLogger.Wrap(err)
		}
	} else {
		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(zap.InfoLevel)
		config.DisableStacktrace = true
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		var err error
		logger, err = config.Build()
		if err != nil {
			return errcode.ErrInitLogger.Wrap(err)
		}
	}
	if zipkinEndpoint != "" {
		reporter := reporterhttp.NewReporter(zipkinEndpoint)
		localEndpoint := &model.Endpoint{ServiceName: "pathwar"}
		sampler, err := zipkin.NewCountingSampler(1)
		if err != nil {
			return errcode.ErrInitTracer.Wrap(err)
		}
		nativeTracer, err := zipkin.NewTracer(
			reporter,
			zipkin.WithSampler(sampler),
			zipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			return errcode.ErrInitTracer.Wrap(err)
		}
		tracer = zipkinot.Wrap(nativeTracer)
		opentracing.SetGlobalTracer(tracer)
	}
	return nil
}

func httpClientFromEnv(ctx context.Context) (*pwapi.HTTPClient, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Timeout: 5 * time.Second})

	conf := &oauth2.Config{
		ClientID:     ssoOpts.ClientID,
		ClientSecret: ssoOpts.ClientSecret,
		Scopes:       []string{"email", "offline_access", "profile", "roles"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  pwsso.KeycloakBaseURL + "/auth/realms/" + ssoOpts.Realm + "/protocol/openid-connect/auth",
			TokenURL: pwsso.KeycloakBaseURL + "/auth/realms/" + ssoOpts.Realm + "/protocol/openid-connect/token",
		},
	}

	if _, err := os.Stat(ssoOpts.TokenFile); err != nil {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n\nthen, write the code in the terminal.\n\n", url)
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			return nil, err
		}

		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			return nil, err
		}

		jsonText, err := json.Marshal(tok)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(ssoOpts.TokenFile, jsonText, 0777); err != nil {
			return nil, err
		}
	}

	byt, err := ioutil.ReadFile(ssoOpts.TokenFile)
	if err != nil {
		return nil, err
	}
	tok := new(oauth2.Token)
	if err = json.Unmarshal(byt, tok); err != nil {
		return nil, err
	}
	ts := conf.TokenSource(ctx, tok)
	_, err = ts.Token()
	if err != nil {
		return nil, err
	}

	return pwapi.NewHTTPClient(oauth2.NewClient(ctx, ts), httpAPIAddr), nil
}

func initSentryFromEnv(startMessage string) (func(), error) {
	cleanup := func() {}
	if globalSentryDSN != "" {
		// doc here: https://docs.sentry.io/platforms/go/config/
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              globalSentryDSN,
			Release:          pwversion.Version,
			AttachStacktrace: true,
			Debug:            false,
			// List of regexp strings that will be used to match against event's message
			// and if applicable, caught errors type and value.
			// If the match is found, then a whole event will be dropped.
			// IgnoreErrors: []string{},
		})
		if err != nil {
			return nil, err
		}
		cleanup = func() {
			sentry.Flush(2 * time.Second)
			sentry.Recover()
		}
		if startMessage != "" {
			sentry.CaptureMessage(startMessage)
		}
	}
	return cleanup, nil
}
