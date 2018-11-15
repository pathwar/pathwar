package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serverOptions struct {
	GRPCBind       string
	HTTPBind       string
	WithReflection bool
}

func (opts serverOptions) String() string {
	out, _ := json.Marshal(opts)
	return string(out)
}

func server(opts *serverOptions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errs := make(chan error)
	go func() { errs <- errors.Wrap(startGRPCServer(ctx, opts), "gRPC server error") }()
	go func() { errs <- errors.Wrap(startHTTPServer(ctx, opts), "HTTP server error") }()
	return <-errs
}

func startHTTPServer(ctx context.Context, opts *serverOptions) error {
	mux := runtime.NewServeMux()
	grpcOpts := []grpc.DialOption{grpc.WithInsecure()}
	if err := RegisterServerHandlerFromEndpoint(ctx, mux, opts.GRPCBind, grpcOpts); err != nil {
		return err
	}
	zap.L().Info("starting HTTP server", zap.String("bind", opts.HTTPBind))
	return http.ListenAndServe(opts.HTTPBind, mux)
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
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(grpcLogger),
		grpc_recovery.StreamServerInterceptor(),
	}
	serverUnaryOpts := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(grpcLogger),
		grpc_recovery.UnaryServerInterceptor(),
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(serverStreamOpts...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(serverUnaryOpts...)),
	)
	RegisterServerServer(grpcServer, &svc{})
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
