package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func serverSetupFlags(flags *pflag.FlagSet, opts *Options) {
	flags.StringVar(&opts.GRPCBind, "grpc-bind", ":9111", "gRPC server address")
	flags.StringVar(&opts.HTTPBind, "http-bind", ":8000", "HTTP server address")
	flags.StringVar(&opts.JWTKey, "jwt-key", "", "JWT secure key")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func NewServerCommand() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}
			return server(opts)
		},
	}
	serverSetupFlags(cmd.Flags(), opts)
	return cmd
}
