package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"net"
	"net/http"

	"github.com/gogo/gateway"
	"github.com/gogo/protobuf/gogoproto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"pathwar.pw/sql"
)

var _ = gogoproto.IsStdTime

type Options struct {
	GRPCBind       string
	HTTPBind       string
	JWTKey         string
	WithReflection bool
}

func (opts Options) String() string {
	out, _ := json.Marshal(opts)
	return string(out)
}

func server(opts *Options) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errs := make(chan error)
	go func() { errs <- errors.Wrap(startGRPCServer(ctx, opts), "gRPC server error") }()
	go func() { errs <- errors.Wrap(startHTTPServer(ctx, opts), "HTTP server error") }()
	return <-errs
}

func startHTTPServer(ctx context.Context, opts *Options) error {
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
	zap.L().Info("starting HTTP server", zap.String("bind", opts.HTTPBind))
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	return http.ListenAndServe(opts.HTTPBind, mux)
}

func startGRPCServer(ctx context.Context, opts *Options) error {
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

	db, err := sql.FromOpts(sql.GetOptions())
	if err != nil {
		return errors.Wrap(err, "failed to initialize database")
	}

	svc, err := newSvc(opts, db)
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

func newSvc(opts *Options, db *gorm.DB) (*svc, error) {
	jwtKey := []byte(opts.JWTKey)
	if len(jwtKey) == 0 { // generate random JWT key
		jwtKey = make([]byte, 128)
		if _, err := rand.Read(jwtKey); err != nil {
			return nil, errors.Wrap(err, "failed to generate random JWT token")
		}
	}
	return &svc{
		jwtKey: jwtKey,
		db:     db,
	}, nil
}
