package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func serverSetupFlags(flags *pflag.FlagSet, opts *serverOptions) {
	flags.StringVar(&opts.GRPCBind, "grpc-bind", ":9111", "gRPC server address")
	flags.StringVar(&opts.HTTPBind, "http-bind", ":8000", "HTTP server address")
	viper.BindPFlags(flags)
}

func NewServerCommand() *cobra.Command {
	opts := &serverOptions{}
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
