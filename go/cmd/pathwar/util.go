package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Bearer/bearer-go"
	"github.com/adrg/xdg"
	"github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql" // required by gorm
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"moul.io/srand"
	"moul.io/zapconfig"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwagent"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
	"pathwar.land/pathwar/v2/go/pkg/pwversion"
)

const (
	defaultDBURN       = "root:uns3cur3@tcp(127.0.0.1:3306)/pathwar?charset=utf8mb4&parseTime=true"
	defaultHTTPApiAddr = "https://api.pathwar.land"
)

var (
	logger *zap.Logger
	tracer opentracing.Tracer

	flagOutput = os.Stderr

	// flag vars
	agentOpts  = pwagent.NewOpts()
	serverOpts = pwapi.NewServerOpts()
	ssoOpts    = pwsso.NewOpts()

	DBURN           string
	DBMaxOpenTries  int
	bearerSecretKey string
	composePSDepth  int
	globalDebug     bool
	globalSentryDSN string
	httpAPIAddr     string
	zipkinEndpoint  string
)

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
	config := zapconfig.Configurator{}
	if globalDebug {
		config.SetLevel(zap.DebugLevel)
	} else {
		config.SetLevel(zap.InfoLevel)
	}
	var err error
	logger, err = config.Build()
	if err != nil {
		return errcode.ErrInitLogger.Wrap(err)
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

	dir, file := filepath.Split(ssoOpts.TokenFile)
	if dir == "" {
		dir = filepath.Join(xdg.ConfigHome, "pathwar", "tokens")
	}

	ssoOpts.TokenFile = filepath.Join(dir, file)

	if _, err := os.Stat(ssoOpts.TokenFile); err != nil {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}

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
