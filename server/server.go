package server

import (
	"context"
	"encoding/json"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type serverOptions struct {
	GRPCBind string
}

func (opts serverOptions) String() string {
	out, _ := json.Marshal(opts)
	return string(out)
}

func server(opts *serverOptions) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := RegisterServerFromEndpoint(ctx, mux, opts.HttpBind, opts)
	listener, err := net.Listen("tcp", opts.GRPCBind)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	grpcServer := grpc.NewServer()
	RegisterServerServer(grpcServer, &svc{})

	zap.L().Info("grpc server started", zap.String("bind", opts.GRPCBind))
	grpcServer.Serve(listener)
	return nil
}
