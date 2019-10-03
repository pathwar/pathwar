package pwserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gogo/gateway"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oklog/run"
	"github.com/rs/cors"
	chilogger "github.com/treastech/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"pathwar.land/go/pkg/pwengine"
)

type Opts struct {
	Logger             *zap.Logger
	GRPCBind           string
	HTTPBind           string
	CORSAllowedOrigins string
	RequestTimeout     time.Duration
	ShutdownTimeout    time.Duration
}

func Start(ctx context.Context, engine pwengine.Engine, opts Opts) (func() error, func(), error) {
	// assign default opts
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.CORSAllowedOrigins == "" {
		opts.CORSAllowedOrigins = "*"
	}
	if opts.GRPCBind == "" {
		opts.GRPCBind = ":9111" // FIXME: get random port
	}
	if opts.HTTPBind == "" {
		opts.HTTPBind = ":8000" // FIXME: get random port
	}
	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = 5 * time.Second
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 6 * time.Second
	}

	var (
		g          run.Group
		grpcLogger = opts.Logger.Named("grpc")
		httpLogger = opts.Logger.Named("http")
	)

	grpcListener, err := net.Listen("tcp", opts.GRPCBind)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start gRPC listener: %w", err)
	}
	{ // gRPC server
		authFunc := func(context.Context) (context.Context, error) {
			// we use svc.AuthFuncOverride to manage authentication
			//
			// this code should never be reached
			return nil, pwengine.ErrNotImplemented
		}
		serverStreamOpts := []grpc.StreamServerInterceptor{
			grpc_recovery.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(authFunc),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(grpcLogger),
			grpc_recovery.StreamServerInterceptor(),
		}
		serverUnaryOpts := []grpc.UnaryServerInterceptor{
			grpc_recovery.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(authFunc),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(grpcLogger),
			grpc_recovery.UnaryServerInterceptor(),
		}
		grpcServer := grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(serverStreamOpts...)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(serverUnaryOpts...)),
		)
		pwengine.RegisterEngineServer(grpcServer, engine)
		g.Add(func() error {
			grpcLogger.Debug("starting gRPC server", zap.String("bind", opts.GRPCBind))
			return grpcServer.Serve(grpcListener)
		}, func(error) {
			grpcServer.GracefulStop()
		})
	}
	{ // HTTP server
		r := chi.NewRouter()
		cors := cors.New(cors.Options{
			AllowedOrigins:   strings.Split(opts.CORSAllowedOrigins, ","),
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		r.Use(cors.Handler)
		r.Use(chilogger.Logger(httpLogger))
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(opts.RequestTimeout))
		r.Use(middleware.RealIP)
		r.Use(middleware.RequestID)
		gwmux := runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.JSONPb{
				EmitDefaults: false,
				Indent:       "  ",
				OrigName:     true,
			}),
			runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
		)
		grpcOpts := []grpc.DialOption{grpc.WithInsecure()}
		if err := pwengine.RegisterEngineHandlerFromEndpoint(ctx, gwmux, opts.GRPCBind, grpcOpts); err != nil {
			return nil, nil, fmt.Errorf("failed to register service on gateway: %w", err)
		}
		r.Mount("/", gwmux)
		srv := http.Server{
			Addr:    opts.HTTPBind,
			Handler: r,
		}
		g.Add(func() error {
			httpLogger.Debug("starting HTTP server", zap.String("bind", opts.HTTPBind))
			return srv.ListenAndServe()
		}, func(error) {
			ctx, cancel := context.WithTimeout(ctx, opts.ShutdownTimeout)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				httpLogger.Warn("failed to shutdown HTTP server", zap.Error(err))
			}
		})
	}

	cleaner := func() {
		if err := grpcListener.Close(); err != nil {
			grpcLogger.Warn("failed to close gRPC listener", zap.Error(err))
		}
	}

	return g.Run, cleaner, nil
}
