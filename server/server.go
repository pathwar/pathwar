package server // import "pathwar.land/server"

import (
	"context"
	"crypto/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gogo/gateway"
	"github.com/gogo/protobuf/gogoproto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"pathwar.land/sql"
)

var _ = gogoproto.IsStdTime

func server(opts *serverOptions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errs := make(chan error)
	go func() { errs <- errors.Wrap(startGRPCServer(ctx, opts), "gRPC server error") }()
	go func() { errs <- errors.Wrap(startHTTPServer(ctx, opts), "HTTP server error") }()
	return <-errs
}

func startHTTPServer(ctx context.Context, opts *serverOptions) error {
	// configure gateway mux
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.JSONPb{
			EmitDefaults: false,
			Indent:       "  ",
			OrigName:     true,
		}),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	grpcOpts := []grpc.DialOption{grpc.WithInsecure()}
	if err := RegisterServerHandlerFromEndpoint(ctx, gwmux, opts.GRPCBind, grpcOpts); err != nil {
		return err
	}

	// configure chi router
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // FIXME: if production, should be the production portal
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	// gateway mux
	r.Mount(opts.APIPrefix, http.StripPrefix(strings.TrimRight(opts.APIPrefix, "/"), gwmux))
	// static files
	if opts.WebDir != "" {
		fs := http.StripPrefix("/", http.FileServer(http.Dir(opts.WebDir)))
		r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}))
	}

	// start HTTP server
	zap.L().Info("starting HTTP server", zap.String("bind", opts.HTTPBind))
	return http.ListenAndServe(opts.HTTPBind, r)
}

func startGRPCServer(ctx context.Context, opts *serverOptions) error {
	listener, err := net.Listen("tcp", opts.GRPCBind)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	defer func() {
		if err := listener.Close(); err != nil {
			zap.L().Error(
				"failed to close listener",
				zap.String("address", opts.GRPCBind),
				zap.Error(err),
			)
		}
	}()

	grpcLogger := zap.L().Named("grpc")
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

	svc, err := newSvc(opts)
	if err != nil {
		return errors.Wrap(err, "failed to initialize service")
	}
	RegisterServerServer(grpcServer, svc)
	if opts.WithReflection {
		reflection.Register(grpcServer)
	}

	go func() {
		defer grpcServer.GracefulStop()
		<-ctx.Done()
	}()

	zap.L().Info("starting gRPC server", zap.String("bind", opts.GRPCBind))
	return grpcServer.Serve(listener)
}

func newSvc(opts *serverOptions) (*svc, error) {
	db, err := sql.FromOpts(&opts.sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize database")
	}

	jwtKey := []byte(opts.JWTKey)
	if len(jwtKey) == 0 { // generate random JWT key
		jwtKey = make([]byte, 128)
		if _, err := rand.Read(jwtKey); err != nil {
			return nil, errors.Wrap(err, "failed to generate random JWT token")
		}
	}
	return &svc{
		jwtKey:    jwtKey,
		db:        db,
		startedAt: time.Now(),
	}, nil
}
