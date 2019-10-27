package pwserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
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
	WithPprof          bool
}

type Server struct {
	grpcServer       *grpc.Server
	grpcListenerAddr string
	httpListenerAddr string
	logger           *zap.Logger
	workers          run.Group
}

func New(ctx context.Context, engine pwengine.Engine, opts Opts) (*Server, error) {
	// assign default opts
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.CORSAllowedOrigins == "" {
		opts.CORSAllowedOrigins = "*"
	}
	if opts.GRPCBind == "" {
		opts.GRPCBind = ""
	}
	if opts.HTTPBind == "" {
		opts.HTTPBind = ":0"
	}
	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = 5 * time.Second
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 6 * time.Second
	}

	var (
		grpcLogger = opts.Logger.Named("grpc")
		httpLogger = opts.Logger.Named("http")
		server     = Server{
			logger: opts.Logger,
		}
	)

	{ // local gRPC server
		authFunc := func(context.Context) (context.Context, error) { return nil, pwengine.ErrNotImplemented }
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
		server.grpcServer = grpcServer
	}

	if opts.HTTPBind != "" || opts.GRPCBind != "" { // grpcbind is required for grpc-gateway (for now)
		grpcListener, err := net.Listen("tcp", opts.GRPCBind)
		if err != nil {
			return nil, fmt.Errorf("start gRPC listener: %w", err)
		}
		server.grpcListenerAddr = grpcListener.Addr().String()

		server.workers.Add(func() error {
			grpcLogger.Debug("starting gRPC server", zap.String("bind", opts.GRPCBind))
			return server.grpcServer.Serve(grpcListener)
		}, func(error) {
			if err := grpcListener.Close(); err != nil {
				grpcLogger.Warn("close gRPC listener", zap.Error(err))
			}
		})
	}

	if opts.HTTPBind != "" {
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
		if err := pwengine.RegisterEngineHandlerFromEndpoint(ctx, gwmux, server.grpcListenerAddr, grpcOpts); err != nil {
			return nil, fmt.Errorf("register service on gateway: %w", err)
		}
		r.Mount("/", gwmux)
		if opts.WithPprof {
			r.HandleFunc("/debug/pprof/*", pprof.Index)
			r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			r.HandleFunc("/debug/pprof/profile", pprof.Profile)
			r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		}
		http.DefaultServeMux = http.NewServeMux() // disables default handlers registere by importing net/http/pprof for security reasons
		listener, err := net.Listen("tcp", opts.HTTPBind)
		if err != nil {
			return nil, fmt.Errorf("start HTTP listener: %w", err)
		}
		server.httpListenerAddr = listener.Addr().String()
		srv := http.Server{
			Handler: r,
		}
		server.workers.Add(func() error {
			httpLogger.Debug("starting HTTP server", zap.String("bind", opts.HTTPBind))
			return srv.Serve(listener)
		}, func(error) {
			ctx, cancel := context.WithTimeout(ctx, opts.ShutdownTimeout)
			if err := srv.Shutdown(ctx); err != nil {
				httpLogger.Warn("shutdown HTTP server", zap.Error(err))
			}
			defer cancel()
			if err := listener.Close(); err != nil {
				httpLogger.Warn("close HTTP listener", zap.Error(err))
			}
		})
	}

	// FIXME: add gRPC web support

	return &server, nil
}

func (s *Server) Run() error {
	return s.workers.Run()
}

func (s *Server) Close() {
	s.grpcServer.GracefulStop()
}

func (s *Server) HTTPListenerAddr() string { return s.httpListenerAddr }
func (s *Server) GRPCListenerAddr() string { return s.grpcListenerAddr }
