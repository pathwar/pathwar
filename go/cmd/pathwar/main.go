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
	"pathwar.land/go/pkg/pwchallenge"
	"pathwar.land/go/pkg/pwcompose"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwengine"
	"pathwar.land/go/pkg/pwserver"
	"pathwar.land/go/pkg/pwsso"
	"pathwar.land/go/pkg/pwversion"
)

const (
	defaultSSOPubKey    = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlEFxLlywsbI5BQ7DVkA66fICWGIYPpD+aZNYRR7SIc0zdtJR4xMOt5CjM0vbYT4z2a1U2yl0ewunyxFm8niS8w6mKYFnOS4nnSchQyIAmJkpLC4eAjijCdEHdr8mSqamThSrVRGSYEEsa+adidC13kRDy7NDKhvZb8F0YqnktNk6WHSlb8r2QRLPJ1DX534jjXPY6l/eoHuLJAOZxBlfwV5Dg37TVmf2xAH812E7ZigycLAvhsMvr5x2jLavAEEnZZmlQf4cyQ4tlMzKS1Zp0NcdOGS/i6lrndc5pNtZQuGr8IGBrEbTRFUiavn/HDnyalYZy8T5LakXRdVaKdshAQIDAQAB"
	defaultSSORealm     = "Pathwar-Dev"
	defaultSSOClientID  = "platform-cli"
	defaultDBURN        = "root:uns3cur3@tcp(127.0.0.1:3306)/pathwar?charset=utf8&parseTime=true"
	defaultDockerPrefix = "pathwar/"
)

var (
	logger *zap.Logger
	// global flags
	globalFlags = flag.NewFlagSet("pathwar", flag.ExitOnError)
	globalDebug = globalFlags.Bool("debug", false, "debug mode")

	// engine flags
	engineFlags          = flag.NewFlagSet("engine", flag.ExitOnError)
	engineDBURN          = engineFlags.String("urn", defaultDBURN, "MySQL URN")
	engineSSOPubkey      = engineFlags.String("sso-pubkey", "", "SSO Public Key")
	engineSSORealm       = engineFlags.String("sso-realm", defaultSSORealm, "SSO Realm")
	engineSSOClientID    = engineFlags.String("sso-clientid", defaultSSOClientID, "SSO ClientID")
	engineSSOAllowUnsafe = engineFlags.Bool("sso-unsafe", false, "Allow unsafe SSO")

	// sso flags
	ssoFlags       = flag.NewFlagSet("sso", flag.ExitOnError)
	ssoPubkey      = ssoFlags.String("pubkey", "", "SSO Public Key")
	ssoRealm       = ssoFlags.String("realm", defaultSSORealm, "SSO Realm")
	ssoClientID    = ssoFlags.String("clientid", defaultSSOClientID, "SSO ClientID")
	ssoAllowUnsafe = ssoFlags.Bool("unsafe", false, "Allow unsafe SSO")

	// compose flags
	composeFlags = flag.NewFlagSet("compose", flag.ExitOnError)

	// compose prepare flags
	composePrepareFlags   = flag.NewFlagSet("compose prepare", flag.ExitOnError)
	composePrepareNoPush  = composePrepareFlags.Bool("no-push", false, "don't push images")
	composePreparePrefix  = composePrepareFlags.String("prefix", defaultDockerPrefix, "docker image prefix")
	composePrepareVersion = composePrepareFlags.String("version", "1.0.0", "challenge version")

	// compose up flags
	composeUpFlags       = flag.NewFlagSet("compose up", flag.ExitOnError)
	composeUpInstanceKey = composeUpFlags.String("instance-key", "default", "instance key used to generate instance ID")

	// compose down flags
	composeDownFlags        = flag.NewFlagSet("compose down", flag.ExitOnError)
	composeDownRemoveImages = composePrepareFlags.Bool("rmi", false, "remove images as well")
	composeDownKeepVolumes  = composePrepareFlags.Bool("keep-volumes", false, "keep volumes")

	// compose ps flags
	composePSFlags = flag.NewFlagSet("compose ps", flag.ExitOnError)
	composePSDepth = composePSFlags.Int("depth", 0, "depth to display")

	// hypervisor flags
	hypervisorFlags = flag.NewFlagSet("hypervisor", flag.ExitOnError)

	// server flags
	serverFlags              = flag.NewFlagSet("server", flag.ExitOnError)
	serverHTTPBind           = serverFlags.String("http-bind", ":8000", "HTTP server address")
	serverGRPCBind           = serverFlags.String("grpc-bind", ":9111", "gRPC server address")
	serverRequestTimeout     = serverFlags.Duration("request-timeout", 5*time.Second, "request timeout")
	serverShutdownTimeout    = serverFlags.Duration("shutdown-timeout", 6*time.Second, "shutdown timeout")
	serverCORSAllowedOrigins = serverFlags.String("cors-allowed-origins", "*", "allowed CORS origins")
	serverWithPprof          = serverFlags.Bool("with-pprof", false, "enable pprof endpoints")

	// misc flags
	miscFlags = flag.NewFlagSet("misc", flag.ExitOnError)
)

func main() {
	log.SetFlags(0)

	defer func() {
		if logger != nil {
			_ = logger.Sync()
		}
	}()

	globalPreRun := func() error {
		rand.Seed(time.Now().UnixNano())
		if *globalDebug {
			config := zap.NewDevelopmentConfig()
			config.Level.SetLevel(zap.DebugLevel)
			config.DisableStacktrace = true
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			var err error
			logger, err = config.Build()
			if err != nil {
				return fmt.Errorf("init logger: %w", err)
			}
		} else {
			config := zap.NewDevelopmentConfig()
			config.Level.SetLevel(zap.InfoLevel)
			config.DisableStacktrace = true
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			var err error
			logger, err = config.Build()
			if err != nil {
				return fmt.Errorf("init logger: %w", err)
			}
		}
		return nil
	}

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

			// init engine
			engine, _, _, closer, err := engineFromFlags()
			if err != nil {
				return fmt.Errorf("start engine: %w", err)
			}
			defer closer()

			var (
				ctx    = context.Background()
				g      run.Group
				server *pwserver.Server
			)
			{ // server
				opts := pwserver.Opts{
					Logger:             logger.Named("server"),
					GRPCBind:           *serverGRPCBind,
					HTTPBind:           *serverHTTPBind,
					CORSAllowedOrigins: *serverCORSAllowedOrigins,
					RequestTimeout:     *serverRequestTimeout,
					ShutdownTimeout:    *serverShutdownTimeout,
					WithPprof:          *serverWithPprof,
				}
				var err error
				server, err = pwserver.New(ctx, engine, opts)
				if err != nil {
					return fmt.Errorf("init server: %w", err)
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
				zap.String("http-bind", server.HTTPListenerAddr()),
				zap.String("grpc-bind", server.GRPCListenerAddr()),
			)

			if err := g.Run(); err != nil {
				return fmt.Errorf("the group was terminated with: %w", err)
			}
			return nil
		},
	}

	sqldump := &ffcli.Command{
		Name:  "sqldump",
		Usage: "pathwar [global flags] engine [engine flags] sqldump",
		Exec: func([]string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			// init engine
			_, db, _, closer, err := engineFromFlags()
			if err != nil {
				return fmt.Errorf("start engine: %w", err)
			}
			defer closer()

			dump, err := pwdb.GetDump(db)
			if err != nil {
				return fmt.Errorf("dump database: %w", err)
			}

			out, _ := json.MarshalIndent(dump, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}

	sqlinfo := &ffcli.Command{
		Name:  "sqlinfo",
		Usage: "pathwar [global flags] engine [engine flags] sqlinfo",
		Exec: func([]string) error {
			if err := globalPreRun(); err != nil {
				return err
			}

			// init engine
			_, db, _, closer, err := engineFromFlags()
			if err != nil {
				return fmt.Errorf("start engine: %w", err)
			}
			defer closer()

			info, err := pwdb.GetInfo(db, logger)
			if err != nil {
				return fmt.Errorf("get database info: %w", err)
			}

			out, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}

	engine := &ffcli.Command{
		Name:        "engine",
		Usage:       "pathwar [global flags] engine [engine flags] <subcommand> [flags] [args...]",
		ShortHelp:   "manage the Pathwar engine",
		FlagSet:     engineFlags,
		Subcommands: []*ffcli.Command{server, sqldump, sqlinfo},
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	pwchallengeBinary := &ffcli.Command{
		Name:  "pwchallenge-binary",
		Usage: "pathwar [global flags] misc [misc flags] pwchallenge-binary",
		Exec: func([]string) error {
			binary, err := pwchallenge.Binary()
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
		Subcommands: []*ffcli.Command{pwchallengeBinary},
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
				return fmt.Errorf("get sso client from flags: %w", err)
			}

			// whoami
			info, err := sso.Whoami(args[0])
			if err != nil {
				return fmt.Errorf("get 'whoami' from SSO: %w", err)
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
				return fmt.Errorf("get sso client from flags: %w", err)
			}

			// logout
			if err := sso.Logout(args[0]); err != nil {
				return fmt.Errorf("logout from SSO: %w", err)
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
				return fmt.Errorf("get sso client from flags: %w", err)
			}

			// token
			token, _, err := sso.TokenWithClaims(args[0])
			if err != nil {
				return fmt.Errorf("get claims: %w", err)
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
			return pwcompose.Prepare(
				path,
				*composePreparePrefix,
				*composePrepareNoPush,
				*composePrepareVersion,
				logger,
			)
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

			return pwcompose.Up(string(preparedCompose), *composeUpInstanceKey, logger)
		},
	}

	composeDown := &ffcli.Command{
		Name:    "down",
		Usage:   "pathwar [global flags] compose [compose flags] down [flags] ID [ID...]",
		FlagSet: composeDownFlags,
		Exec: func(args []string) error {
			if len(args) < 1 {
				return flag.ErrHelp
			}

			if err := globalPreRun(); err != nil {
				return err
			}

			ctx := context.Background()
			cli, err := client.NewEnvClient()
			if err != nil {
				return fmt.Errorf("docker client: %w", err)
			}

			return pwcompose.Down(
				ctx,
				args,
				*composeDownRemoveImages,
				!*composeDownKeepVolumes,
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
				return fmt.Errorf("docker client: %w", err)
			}

			return pwcompose.PS(ctx, *composePSDepth, cli, logger)
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

	hypervisor := &ffcli.Command{
		Name:        "hypervisor",
		Usage:       "pathwar [global flags] hypervisor [sso flags] <subcommand> [flags] [args...]",
		ShortHelp:   "manage an hypervisor node (multiple challenges)",
		Subcommands: []*ffcli.Command{},
		FlagSet:     hypervisorFlags,
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	root := &ffcli.Command{
		Usage:       "pathwar [global flags] <subcommand> [flags] [args...]",
		FlagSet:     globalFlags,
		LongHelp:    "More info here: https://github.com/pathwar/pathwar/wiki/CLI",
		Options:     []ff.Option{ff.WithEnvVarPrefix("PATHWAR")},
		Subcommands: []*ffcli.Command{engine, compose, hypervisor, sso, misc, version},
		Exec:        func([]string) error { return flag.ErrHelp },
	}

	if err := root.Run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatalf("fatal: %+v", err)
	}
}

func engineFromFlags() (pwengine.Engine, *gorm.DB, pwsso.Client, func(), error) {
	// init database
	db, err := gorm.Open("mysql", *engineDBURN)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("init database: %w", err)
	}
	sfn, err := snowflake.NewNode(1)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("init snowflake: %w", err)
	}
	dbOpts := pwdb.Opts{
		Logger: logger.Named("gorm"),
	}
	db, err = pwdb.Configure(db, sfn, dbOpts)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("configure database: %w", err)
	}

	// init SSO
	ssoOpts := pwsso.Opts{
		AllowUnsafe: *engineSSOAllowUnsafe,
		Logger:      logger.Named("sso"),
		ClientID:    *engineSSOClientID,
	}
	if *engineSSOPubkey == "" {
		*engineSSOPubkey = defaultSSOPubKey
	}
	sso, err := pwsso.New(*engineSSOPubkey, *engineSSORealm, ssoOpts)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("init SSO client: %w", err)
	}

	// init engine
	engineOpts := pwengine.Opts{
		Logger: logger.Named("engine"),
	}

	engine, err := pwengine.New(db, sso, engineOpts)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("init engine: %w", err)
	}

	closeFunc := func() {
		if err := engine.Close(); err != nil {
			logger.Warn("close engine", zap.Error(err))
		}
		if err := db.Close(); err != nil {
			logger.Warn("closed database", zap.Error(err))
		}
	}

	logger.Debug("engine initd", zap.Any("db", db), zap.Any("opts", engineOpts))
	return engine, db, sso, closeFunc, nil
}

func ssoFromFlags() (pwsso.Client, error) {
	ssoOpts := pwsso.Opts{
		AllowUnsafe: *ssoAllowUnsafe,
		Logger:      logger.Named("sso"),
		ClientID:    *ssoClientID,
	}
	if *ssoPubkey == "" {
		*ssoPubkey = defaultSSOPubKey
	}
	sso, err := pwsso.New(*ssoPubkey, *ssoRealm, ssoOpts)
	if err != nil {
		return nil, fmt.Errorf("init SSO client: %w", err)
	}
	return sso, nil
}
