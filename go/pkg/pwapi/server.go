package pwapi

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
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oklog/run"
	"github.com/rs/cors"
	"github.com/soheilhy/cmux"
	chilogger "github.com/treastech/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"pathwar.land/go/v2/pkg/errcode"
)

func NewServer(ctx context.Context, svc Service, opts ServerOpts) (*Server, error) {
	// assign default opts
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.CORSAllowedOrigins == "" {
		opts.CORSAllowedOrigins = "*"
	}
	if opts.Bind == "" {
		opts.Bind = ":0"
	}
	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = 5 * time.Second
	}
	if opts.ShutdownTimeout == 0 {
		opts.ShutdownTimeout = 6 * time.Second
	}
	s := Server{logger: opts.Logger}

	// listener
	var err error
	s.masterListener, err = net.Listen("tcp", opts.Bind)
	if err != nil {
		return nil, errcode.ErrServerListen.Wrap(err)
	}
	s.cmux = cmux.New(s.masterListener)
	s.grpcListener = s.cmux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	s.httpListener = s.cmux.Match(cmux.HTTP2(), cmux.HTTP1())
	// FIXME: add gRPC web support
	// FIXME: websocket

	// grpc server
	s.grpcServer, err = grpcServer(svc, opts)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	s.workers.Add(func() error {
		err := s.grpcServer.Serve(s.grpcListener)
		if err != cmux.ErrListenerClosed {
			return err
		}
		return nil
	}, func(error) {
		if err := s.grpcListener.Close(); err != nil {
			opts.Logger.Warn("close listener", zap.Error(err))
		}
	})

	// http server
	httpServer, err := httpServer(ctx, s.ListenerAddr(), opts)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	s.workers.Add(func() error {
		err := httpServer.Serve(s.httpListener)
		if err != cmux.ErrListenerClosed {
			return err
		}
		return nil
	}, func(error) {

		ctx, cancel := context.WithTimeout(ctx, opts.ShutdownTimeout)
		if err := httpServer.Shutdown(ctx); err != nil {
			opts.Logger.Warn("shutdown HTTP server", zap.Error(err))
		}
		defer cancel()
		if err := s.httpListener.Close(); err != nil {
			opts.Logger.Warn("close listener", zap.Error(err))
		}
	})

	s.cmux.HandleError(func(err error) bool {
		s.logger.Warn("cmux error", zap.Error(err))
		return true
	})

	// mux
	s.workers.Add(
		func() error {
			err := s.cmux.Serve()
			return err
		},
		func(err error) {
			fmt.Println(err)
		},
	)
	return &s, nil
}

// Server is an HTTP+gRPC frontend for Service
type Server struct {
	grpcServer     *grpc.Server
	masterListener net.Listener
	grpcListener   net.Listener
	httpListener   net.Listener
	cmux           cmux.CMux
	logger         *zap.Logger
	workers        run.Group
}

type ServerOpts struct {
	Logger             *zap.Logger
	Bind               string
	CORSAllowedOrigins string
	RequestTimeout     time.Duration
	ShutdownTimeout    time.Duration
	WithPprof          bool
}

func (s *Server) Run() error {
	return s.workers.Run()
}

func (s *Server) Close() {
	//go s.grpcServer.GracefulStop()
	//time.Sleep(time.Second)
	//s.grpcServer.Stop()
	s.masterListener.Close()
}

func (s *Server) ListenerAddr() string {
	return s.masterListener.Addr().String()
}

func grpcServer(svc Service, opts ServerOpts) (*grpc.Server, error) {
	logger := opts.Logger.Named("grpc")
	authFunc := func(context.Context) (context.Context, error) {
		return nil, errcode.ErrNotImplemented
	}
	serverStreamOpts := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
		grpc_auth.StreamServerInterceptor(authFunc),
		//grpc_ctxtags.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(logger),
		grpc_recovery.StreamServerInterceptor(),
	}
	serverUnaryOpts := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_auth.UnaryServerInterceptor(authFunc),
		//grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_recovery.UnaryServerInterceptor(),
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(serverStreamOpts...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(serverUnaryOpts...)),
	)
	RegisterServiceServer(grpcServer, svc)

	return grpcServer, nil
}

func httpServer(ctx context.Context, serverListenerAddr string, opts ServerOpts) (*http.Server, error) {
	logger := opts.Logger.Named("http")
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
	r.Use(chilogger.Logger(logger))
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
	err := RegisterServiceHandlerFromEndpoint(ctx, gwmux, serverListenerAddr, grpcOpts)
	if err != nil {
		return nil, errcode.ErrServerRegisterGateway.Wrap(err)
	}
	r.Mount("/", gwmux)
	if opts.WithPprof {
		r.HandleFunc("/debug/pprof/*", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
	http.DefaultServeMux = http.NewServeMux() // disables default handlers registered by importing net/http/pprof for security reasons

	return &http.Server{
		Addr:    ":8000",
		Handler: r,
	}, nil
}
