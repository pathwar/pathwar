package server

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type serverOptions struct{}

func (opts serverOptions) String() string {
	out, _ := json.Marshal(opts)
	return string(out)
}

func serverSetupFlags(flags *pflag.FlagSet, opts *serverOptions) {
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

func server(opts *serverOptions) error {
	fmt.Println("hello world server")
	return nil
}
