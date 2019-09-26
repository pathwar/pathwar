package server

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.land/client"
	"pathwar.land/pkg/cli"
	"pathwar.land/sql"
)

type serverOptions struct {
	SQL    sql.Options
	Client client.Options

	GRPCBind       string
	HTTPBind       string
	WithReflection bool
	WebDir         string
	APIPrefix      string
}

func (opts serverOptions) String() string {
	out, _ := json.Marshal(opts)
	return string(out)
}

func Commands() cli.Commands {
	return cli.Commands{
		"server": &serverCommand{},
	}
}

type serverCommand struct{ opts serverOptions }

func (cmd *serverCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *serverCommand) ParseFlags(flags *pflag.FlagSet) {
	flags.StringVar(&cmd.opts.GRPCBind, "grpc-bind", ":9111", "gRPC server address")
	flags.StringVar(&cmd.opts.HTTPBind, "http-bind", ":8000", "HTTP server address")
	flags.StringVar(&cmd.opts.WebDir, "web-dir", "", "Static Files Directory")
	flags.StringVar(&cmd.opts.APIPrefix, "api-prefix", "/api/", "API Prefix")
	flags.BoolVarP(&cmd.opts.WithReflection, "grpc-reflection", "", false, "enable gRPC reflection")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}
func (cmd *serverCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "server",
		RunE: func(_ *cobra.Command, args []string) error {
			opts := cmd.opts
			opts.SQL = sql.GetOptions(commands)
			opts.Client = client.GetOptions(commands)
			return server(&opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	commands["sql"].ParseFlags(cc.Flags())
	commands["client"].ParseFlags(cc.Flags())
	return cc
}
