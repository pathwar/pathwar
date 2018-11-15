package server

import (
	"encoding/json"
	"net"

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
