package main // import "pathwar.land/go/cmd/pathwar"

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
	"os/signal"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/docker/docker/client"
	_ "github.com/go-sql-driver/mysql" // required by gorm
	"github.com/jinzhu/gorm"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwagent"
	"pathwar.land/go/pkg/pwapi"
	"pathwar.land/go/pkg/pwcompose"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwinit"
	"pathwar.land/go/pkg/pwsso"
	"pathwar.land/go/pkg/pwversion"
)

const (
	defaultSSOPubKey       = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlEFxLlywsbI5BQ7DVkA66fICWGIYPpD+aZNYRR7SIc0zdtJR4xMOt5CjM0vbYT4z2a1U2yl0ewunyxFm8niS8w6mKYFnOS4nnSchQyIAmJkpLC4eAjijCdEHdr8mSqamThSrVRGSYEEsa+adidC13kRDy7NDKhvZb8F0YqnktNk6WHSlb8r2QRLPJ1DX534jjXPY6l/eoHuLJAOZxBlfwV5Dg37TVmf2xAH812E7ZigycLAvhsMvr5x2jLavAEEnZZmlQf4cyQ4tlMzKS1Zp0NcdOGS/i6lrndc5pNtZQuGr8IGBrEbTRFUiavn/HDnyalYZy8T5LakXRdVaKdshAQIDAQAB"
	defaultSSORealm        = "Pathwar-Dev"
	defaultSSOClientID     = "platform-cli"
	defaultSSOClientSecret = ""
	defaultDBURN           = "root:uns3cur3@tcp(127.0.0.1:3306)/pathwar?charset=utf8&parseTime=true"
	defaultDockerPrefix    = "pathwar/"
	defaultTokenFile       = "pathwar_oauth_token.json"
	defaultHTTPApiAddr     = "https://api-dev.pathwar.land"
	defaultAgentName       = "localhost"
)

var (
	logger     *zap.Logger
	flagOutput = os.Stderr

	// flag vars
	globalDebug                 bool
	agentForceRecreate          bool
	agentDaemonClean            bool
	agentDaemonRunOnce          bool
	agentDaemonLoopDelay        time.Duration
	agentName                   string
	agentNginxDockerImage       string
	agentNginxDomainSuffix      string
	agentNginxHostIP            string
	agentNginxHostPort          string
	agentNginxModeratorPassword string
	agentNginxSalt              string
	apiDBURN                    string
	composeDownKeepVolumes      bool
	composeDownRemoveImages     bool
	composeDownWithNginx        bool
	composePSDepth              int
	composePrepareNoPush        bool
	composePreparePrefix        string
	composePrepareVersion       string
	composeUpInstanceKey        string
	composeUpForceRecreate      bool
	httpAPIAddr                 string
	serverCORSAllowedOrigins    string
	serverBind                  string
	serverRequestTimeout        time.Duration
	serverShutdownTimeout       time.Duration
	serverWithPprof             bool
	ssoAllowUnsafe              bool
	ssoClientID                 string
	ssoClientSecret             string
	ssoPubkey                   string
	ssoRealm                    string
	ssoTokenFile                string
)

func main() {
	log.SetFlags(0)

	defer func() {
		if logger != nil {
			_ = logger.Sync()
		}
	}()

	// setup flags
	var (
		globalFlags         = flag.NewFlagSet("pathwar", flag.ExitOnError)
		agentDaemonFlags    = flag.NewFlagSet("agent daemon", flag.ExitOnError)
		agentFlags          = flag.NewFlagSet("agent", flag.ExitOnError)
		agentNginxFlags     = flag.NewFlagSet("agent nginx", flag.ExitOnError)
		apiFlags            = flag.NewFlagSet("api", flag.ExitOnError)
		composeDownFlags    = flag.NewFlagSet("compose down", flag.ExitOnError)
		composeFlags        = flag.NewFlagSet("compose", flag.ExitOnError)
		composePSFlags      = flag.NewFlagSet("compose ps", flag.ExitOnError)
		composePrepareFlags = flag.NewFlagSet("compose prepare", flag.ExitOnError)
		composeUpFlags      = flag.NewFlagSet("compose up", flag.ExitOnError)
		miscFlags           = flag.NewFlagSet("misc", flag.ExitOnError)
		serverFlags         = flag.NewFlagSet("server", flag.ExitOnError)
		ssoFlags            = flag.NewFlagSet("sso", flag.ExitOnError)
	)
	globalFlags.SetOutput(flagOutput) // used in main_test.go
	globalFlags.BoolVar(&globalDebug, "debug", false, "debug mode")
	agentDaemonFlags.BoolVar(&agentDaemonClean, "clean", false, "remove all pathwar instances before executing")
	agentDaemonFlags.BoolVar(&agentDaemonRunOnce, "once", false, "run once and don't start daemon loop")
	agentDaemonFlags.DurationVar(&agentDaemonLoopDelay, "delay", 10*time.Second, "delay between each loop iteration")
	agentDaemonFlags.StringVar(&httpAPIAddr, "http-api-addr", defaultHTTPApiAddr, "HTTP API address")
	agentDaemonFlags.StringVar(&ssoClientID, "sso-clientid", defaultSSOClientID, "SSO ClientID")
	agentDaemonFlags.StringVar(&ssoClientSecret, "sso-clientsecret", defaultSSOClientSecret, "SSO ClientSecret")
	agentDaemonFlags.StringVar(&ssoRealm, "sso-realm", defaultSSORealm, "SSO Realm")
	agentDaemonFlags.StringVar(&ssoTokenFile, "sso-token-file", defaultTokenFile, "Token file")
	agentDaemonFlags.StringVar(&agentName, "agent-name", defaultAgentName, "Agent Name")
	agentNginxFlags.StringVar(&agentNginxDockerImage, "docker-image", "docker.io/library/nginx:stable-alpine", "docker image used to generate nginx proxy container")
	agentNginxFlags.StringVar(&agentNginxDomainSuffix, "domain-suffix", "local", "Domain suffix to append")
	agentNginxFlags.StringVar(&agentNginxHostIP, "host", "0.0.0.0", "HTTP listening addr")
	agentNginxFlags.StringVar(&agentNginxHostPort, "port", "8000", "HTTP listening port")
	agentNginxFlags.StringVar(&agentNginxModeratorPassword, "moderator-password", "", "Challenge moderator password")
	agentNginxFlags.StringVar(&agentNginxSalt, "salt", "", "salt used to generate secure hashes (random if empty)")
	apiFlags.BoolVar(&ssoAllowUnsafe, "sso-unsafe", false, "Allow unsafe SSO")
	apiFlags.StringVar(&apiDBURN, "urn", defaultDBURN, "MySQL URN")
	apiFlags.StringVar(&ssoClientID, "sso-clientid", defaultSSOClientID, "SSO ClientID")
	apiFlags.StringVar(&ssoPubkey, "sso-pubkey", "", "SSO Public Key")
	apiFlags.StringVar(&ssoRealm, "sso-realm", defaultSSORealm, "SSO Realm")
	composeDownFlags.BoolVar(&composeDownKeepVolumes, "keep-volumes", false, "keep volumes")
	composeDownFlags.BoolVar(&composeDownRemoveImages, "rmi", false, "remove images as well")
	composeDownFlags.BoolVar(&composeDownWithNginx, "with-nginx", false, "down nginx container and proxy network as well")
	composePSFlags.IntVar(&composePSDepth, "depth", 0, "depth to display")
	composePrepareFlags.BoolVar(&composePrepareNoPush, "no-push", false, "don't push images")
	composePrepareFlags.StringVar(&composePreparePrefix, "prefix", defaultDockerPrefix, "docker image prefix")
	composePrepareFlags.StringVar(&composePrepareVersion, "version", "1.0.0", "challenge version")
	composeUpFlags.StringVar(&composeUpInstanceKey, "instance-key", "default", "instance key used to generate instance ID")
	composeUpFlags.BoolVar(&composeUpForceRecreate, "force-recreate", false, "down previously created instances of challenge")
	serverFlags.BoolVar(&serverWithPprof, "with-pprof", false, "enable pprof endpoints")
	serverFlags.DurationVar(&serverRequestTimeout, "request-timeout", 5*time.Second, "request timeout")
	serverFlags.DurationVar(&serverShutdownTimeout, "shutdown-timeout", 6*time.Second, "shutdown timeout")
	serverFlags.StringVar(&serverCORSAllowedOrigins, "cors-allowed-origins", "*", "allowed CORS origins")
	serverFlags.StringVar(&serverBind, "bind", ":8000", "server address")
	ssoFlags.BoolVar(&ssoAllowUnsafe, "unsafe", false, "Allow unsafe SSO")
	ssoFlags.StringVar(&ssoClientID, "clientid", defaultSSOClientID, "SSO ClientID")
	ssoFlags.StringVar(&ssoPubkey, "pubkey", "", "SSO Public Key")
	ssoFlags.StringVar(&ssoRealm, "realm", defaultSSORealm, "SSO Realm")

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

	server := &ffcli.Command{
		Name:      "server",
		Usage:     "pathwar [global flags] server [server flags] <subcommand> [flags] [args...]",
		ShortHelp: "start a server (HTTP + gRPC)",
		FlagSet:   serverFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			// init svc
			svc, _, _, closer, err := svcFromFlags()
			if err != nil {
				return errcode.ErrStartService.Wrap(err)
			}
			defer closer()

			var (
				ctx    = context.Background()
				g      run.Group
				server *pwapi.Server
			)
			{ // server
				opts := pwapi.ServerOpts{
					Logger:             logger.Named("server"),
					Bind:               serverBind,
					CORSAllowedOrigins: serverCORSAllowedOrigins,
					RequestTimeout:     serverRequestTimeout,
					ShutdownTimeout:    serverShutdownTimeout,
					WithPprof:          serverWithPprof,
				}
				var err error
				server, err = pwapi.NewServer(ctx, svc, opts)
				if err != nil {
					return errcode.ErrInitServer.Wrap(err)
				}
				g.Add(
					server.Run,
					func(error) { server.Close() },
				)
			}
			{ // signal handling and cancellation
				ctx, cancel := context.WithCancel(ctx)
				g.Add(func() error {
					sigch := make(chan os.Signal, 1)
					signal.Notify(sigch, os.Interrupt)
					select {
					case <-sigch:
					case <-ctx.Done():
					}
					return nil
				}, func(error) {
					cancel()
				})
			}

			logger.Info("server started",
				zap.String("bind", server.ListenerAddr()),
			)

			if err := g.Run(); err != nil {
				return errcode.ErrGroupTerminated.Wrap(err)
			}
			return nil
		},
	}

	sqldump := &ffcli.Command{
		Name:  "sqldump",
		Usage: "pathwar [global flags] api [api flags] sqldump",
		Exec: func([]string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			// init svc
			_, db, _, closer, err := svcFromFlags()
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
	}

	sqlinfo := &ffcli.Command{
		Name:  "sqlinfo",
		Usage: "pathwar [global flags] api [api flags] sqlinfo",
		Exec: func([]string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			// init svc
			_, db, _, closer, err := svcFromFlags()
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
	}

	api := &ffcli.Command{
		Name:        "api",
		Usage:       "pathwar [global flags] api [api flags] <subcommand> [flags] [args...]",
		ShortHelp:   "manage the Pathwar API",
		FlagSet:     apiFlags,
		Subcommands: []*ffcli.Command{server, sqldump, sqlinfo},
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	pwinitBinary := &ffcli.Command{
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
	}

	misc := &ffcli.Command{
		Name:        "misc",
		Usage:       "pathwar [global flags] misc [misc flags] <subcommand> [flags] [args...]",
		ShortHelp:   "misc contains advanced commands",
		Subcommands: []*ffcli.Command{pwinitBinary},
		FlagSet:     miscFlags,
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	ssoWhoami := &ffcli.Command{
		Name:  "whoami",
		Usage: "pathwar [global flags] sso [sso flags] whoami TOKEN",
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
	}

	ssoLogout := &ffcli.Command{
		Name:  "logout",
		Usage: "pathwar [global flags] sso [sso flags] logout TOKEN",
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

			// logout
			if err := sso.Logout(args[0]); err != nil {
				return errcode.ErrGetSSOLogout.Wrap(err)
			}
			return nil
		},
	}

	ssoToken := &ffcli.Command{
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
	}

	sso := &ffcli.Command{
		Name:        "sso",
		Usage:       "pathwar [global flags] sso [sso flags] <subcommand> [flags] [args...]",
		ShortHelp:   "manage SSO tokens",
		Subcommands: []*ffcli.Command{ssoLogout, ssoToken, ssoWhoami},
		FlagSet:     ssoFlags,
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	composePrepare := &ffcli.Command{
		Name:    "prepare",
		Usage:   "pathwar [global flags] compose [compose flags] prepare [flags] PATH",
		FlagSet: composePrepareFlags,
		Exec: func(args []string) error {
			if len(args) < 1 {
				return flag.ErrHelp
			}
			path := args[0]
			if err := globalPreRun(); err != nil {
				return err
			}

			preparedComposeData, err := pwcompose.Prepare(
				path,
				composePreparePrefix,
				composePrepareNoPush,
				composePrepareVersion,
				logger,
			)

			fmt.Println(preparedComposeData)

			return err
		},
	}

	composeUp := &ffcli.Command{
		Name:    "up",
		Usage:   "pathwar [global flags] compose [compose flags] up [flags] PATH",
		FlagSet: composeUpFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
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

			return pwcompose.Up(ctx, string(preparedCompose), composeUpInstanceKey, composeUpForceRecreate, nil, cli, logger)
		},
	}

	composeDown := &ffcli.Command{
		Name:    "down",
		Usage:   "pathwar [global flags] compose [compose flags] down [flags] ID [ID...]",
		FlagSet: composeDownFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			cli, err := client.NewEnvClient()
			if err != nil {
				return errcode.ErrInitDockerClient.Wrap(err)
			}

			return pwcompose.Down(
				ctx,
				args,
				composeDownRemoveImages,
				!composeDownKeepVolumes,
				composeDownWithNginx,
				cli,
				logger,
			)
		},
	}

	composePS := &ffcli.Command{
		Name:    "ps",
		Usage:   "pathwar [global flags] compose [compose flags] ps [flags]",
		FlagSet: composePSFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
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

	compose := &ffcli.Command{
		Name:        "compose",
		Usage:       "pathwar [global flags] compose [sso flags] <subcommand> [flags] [args...]",
		Subcommands: []*ffcli.Command{composePrepare, composeUp, composePS, composeDown},
		ShortHelp:   "manage a challenge",
		FlagSet:     composeFlags,
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	agentDaemon := &ffcli.Command{
		Name:    "daemon",
		Usage:   "pathwar [global flags] agent [agent flags] daemon [flags]",
		FlagSet: agentDaemonFlags,
		Exec: func(args []string) error {
			if err := globalPreRun(); err != nil {
				return err
			}
			ctx := context.Background()
			dockerCli, err := client.NewEnvClient()
			if err != nil {
				return errcode.ErrInitDockerClient.Wrap(err)
			}

			apiClient, err := oauthClientFromEnv(ctx)
			if err != nil {
				return errcode.TODO.Wrap(err)
			}

			return pwagent.Daemon(ctx, agentDaemonClean, agentDaemonRunOnce, agentDaemonLoopDelay, dockerCli, apiClient, httpAPIAddr, agentName, logger)
		},
	}

	agentNginx := &ffcli.Command{
		Name:    "nginx",
		Usage:   "pathwar [global flags] agent [agent flags] nginx [flags] ALLOWED_USERS",
		FlagSet: agentNginxFlags,
		Exec: func(args []string) error {
			if len(args) < 1 {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			// prepare AgentOpts
			config := pwagent.AgentOpts{
				HostIP:            agentNginxHostIP,
				HostPort:          agentNginxHostPort,
				DomainSuffix:      agentNginxDomainSuffix,
				ModeratorPassword: agentNginxModeratorPassword,
				Salt:              agentNginxSalt,
				ForceRecreate:     agentForceRecreate,
				NginxDockerImage:  agentNginxDockerImage,
			}
			err := json.Unmarshal([]byte(args[0]), &config.AllowedUsers)
			if err != nil {
				return errcode.ErrInvalidInput.Wrap(err)
			}

			ctx := context.Background()
			cli, err := client.NewEnvClient()
			if err != nil {
				return errcode.ErrInitDockerClient.Wrap(err)
			}

			return pwagent.Nginx(ctx, config, cli, logger)
		},
	}

	agent := &ffcli.Command{
		Name:        "agent",
		Usage:       "pathwar [global flags] agent [sso flags] <subcommand> [flags] [args...]",
		ShortHelp:   "manage an agent node (multiple challenges)",
		Subcommands: []*ffcli.Command{agentDaemon, agentNginx},
		FlagSet:     agentFlags,
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	root := &ffcli.Command{
		Usage:       "pathwar [global flags] <subcommand> [flags] [args...]",
		FlagSet:     globalFlags,
		LongHelp:    "More info here: https://github.com/pathwar/pathwar/wiki/CLI",
		Options:     []ff.Option{ff.WithEnvVarPrefix("PATHWAR")},
		Subcommands: []*ffcli.Command{api, compose, agent, sso, misc, version},
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	if err := root.Run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatalf("fatal: %+v", err)
	}
}

func svcFromFlags() (pwapi.Service, *gorm.DB, pwsso.Client, func(), error) {
	// init database
	db, err := gorm.Open("mysql", apiDBURN)
	if err != nil {
		return nil, nil, nil, nil, errcode.ErrInitDB.Wrap(err)
	}
	sfn, err := snowflake.NewNode(1)
	if err != nil {
		return nil, nil, nil, nil, errcode.ErrInitSnowflake.Wrap(err)
	}
	dbOpts := pwdb.Opts{
		Logger: logger.Named("gorm"),
	}
	db, err = pwdb.Configure(db, sfn, dbOpts)
	if err != nil {
		return nil, nil, nil, nil, errcode.ErrConfigureDB.Wrap(err)
	}

	// init SSO
	sso, err := ssoFromFlags()
	if err != nil {
		return nil, nil, nil, nil, errcode.ErrInitSSOClient.Wrap(err)
	}

	// init svc
	svcOpts := pwapi.ServiceOpts{
		Logger: logger.Named("svc"),
	}

	svc, err := pwapi.NewService(db, sso, svcOpts)
	if err != nil {
		return nil, nil, nil, nil, errcode.ErrInitService.Wrap(err)
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
	return svc, db, sso, closeFunc, nil
}

func ssoFromFlags() (pwsso.Client, error) {
	ssoOpts := pwsso.Opts{
		AllowUnsafe: ssoAllowUnsafe,
		Logger:      logger.Named("sso"),
		ClientID:    ssoClientID,
	}
	if ssoPubkey == "" {
		ssoPubkey = defaultSSOPubKey
	}
	sso, err := pwsso.New(ssoPubkey, ssoRealm, ssoOpts)
	if err != nil {
		return nil, errcode.ErrInitSSOClient.Wrap(err)
	}
	return sso, nil
}

func globalPreRun() error {
	rand.Seed(time.Now().UnixNano())
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
	return nil
}

func oauthClientFromEnv(ctx context.Context) (*http.Client, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Timeout: 5 * time.Second})

	conf := &oauth2.Config{
		ClientID:     ssoClientID,
		ClientSecret: ssoClientSecret,
		Scopes:       []string{"email", "offline_access", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  pwsso.KeycloakBaseURL + "/auth/realms/" + ssoRealm + "/protocol/openid-connect/auth",
			TokenURL: pwsso.KeycloakBaseURL + "/auth/realms/" + ssoRealm + "/protocol/openid-connect/token",
		},
	}

	if _, err := os.Stat(ssoTokenFile); err != nil {
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

		if err := ioutil.WriteFile(ssoTokenFile, jsonText, 0777); err != nil {
			return nil, err
		}
	}

	byt, err := ioutil.ReadFile(ssoTokenFile)
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

	return oauth2.NewClient(ctx, ts), nil
}
