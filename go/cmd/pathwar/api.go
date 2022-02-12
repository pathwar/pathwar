package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"
	"moul.io/banner"
	"moul.io/motd"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func apiCommand() *ffcli.Command {
	var (
		apiFlags    = flag.NewFlagSet("api", flag.ExitOnError)
		serverFlags = flag.NewFlagSet("api server", flag.ExitOnError)
	)
	apiFlags.BoolVar(&ssoOpts.AllowUnsafe, "sso-unsafe", ssoOpts.AllowUnsafe, "Allow unsafe SSO")
	apiFlags.StringVar(&DBURN, "urn", defaultDBURN, "MySQL URN")
	apiFlags.IntVar(&DBMaxOpenTries, "db-max-open-tries", 0, "max DB open tries, unlimited if 0")
	apiFlags.StringVar(&ssoOpts.ClientID, "sso-clientid", ssoOpts.ClientID, "SSO ClientID")
	apiFlags.StringVar(&ssoOpts.Pubkey, "sso-pubkey", ssoOpts.Pubkey, "SSO Public Key")
	apiFlags.StringVar(&ssoOpts.Realm, "sso-realm", ssoOpts.Realm, "SSO Realm")
	serverFlags.BoolVar(&serverOpts.WithPprof, "with-pprof", serverOpts.WithPprof, "enable pprof endpoints")
	serverFlags.DurationVar(&serverOpts.RequestTimeout, "request-timeout", serverOpts.RequestTimeout, "request timeout")
	serverFlags.DurationVar(&serverOpts.ShutdownTimeout, "shutdown-timeout", serverOpts.ShutdownTimeout, "shutdown timeout")
	serverFlags.StringVar(&serverOpts.CORSAllowedOrigins, "cors-allowed-origins", serverOpts.CORSAllowedOrigins, "allowed CORS origins")
	serverFlags.StringVar(&serverOpts.Bind, "bind", serverOpts.Bind, "server address")

	return &ffcli.Command{
		Name:       "api",
		ShortUsage: "pathwar [global flags] api [api flags] <subcommand> [flags] [args...]",
		ShortHelp:  "manage the Pathwar API",
		FlagSet:    apiFlags,
		Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
		Subcommands: []*ffcli.Command{{
			Name:       "server",
			ShortUsage: "pathwar [global flags] server [server flags] <subcommand> [flags] [args...]",
			ShortHelp:  "start a server (HTTP + gRPC)",
			FlagSet:    serverFlags,
			Options:    []ff.Option{ff.WithEnvVarNoPrefix()},
			Exec: func(ctx context.Context, args []string) error {
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
					g      run.Group
					server *pwapi.Server
				)
				g.Add(run.SignalHandler(ctx, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill))
				{ // server
					serverOpts.Tracer = tracer
					serverOpts.Logger = logger.Named("server")
					var err error

					if serverOpts.Bind == "gcloud" {
						serverOpts.Bind = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
						logger.Info("bind", zap.String("address", serverOpts.Bind))
					}

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
			Name:       "sqlinfo",
			ShortUsage: "pathwar [global flags] api [api flags] sqlinfo",
			Exec: func(ctx context.Context, args []string) error {
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
			Name:       "sqldump",
			ShortUsage: "pathwar [global flags] api [api flags] sqldump",
			Exec: func(ctx context.Context, args []string) error {
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
		Exec: func(ctx context.Context, args []string) error { return flag.ErrHelp },
	}
}

func svcFromFlags(logger *zap.Logger) (pwapi.Service, *gorm.DB, func(), error) {
	// init database
	dbConnectTries := 0
dbConnectLoop:
	if DBURN == "gcloud" {
		var (
			dbUser                 = os.Getenv("DB_USER")
			dbPass                 = os.Getenv("DB_PASS")
			instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
			dbName                 = os.Getenv("DB_NAME")
		)
		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}
		DBURN = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPass, socketDir, instanceConnectionName, dbName)
		logger.Info("gcloud URN", zap.String("connection", DBURN))
	}
	db, err := gorm.Open("mysql", DBURN)
	if err != nil {
		dbConnectTries++
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
