package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.land/pkg/cli"
)

type claimsOptions struct{ client Options }
type claimsCommand struct{ opts claimsOptions }

func (cmd *claimsCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "claims",
		Args: func(_ *cobra.Command, args []string) error {
			cmd.opts.client = GetOptions(commands)
			if len(args) == 1 {
				cmd.opts.client.Token = args[0]
			}
			if cmd.opts.client.Token == "" {
				return errors.New("--token is mandatory")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			return runClaims(&cmd.opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	commands["client"].ParseFlags(cc.Flags())
	return cc
}
func (cmd *claimsCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *claimsCommand) ParseFlags(flags *pflag.FlagSet) {
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runClaims(opts *claimsOptions) error {
	token, _, err := TokenWithClaims(opts.client.Token, opts.client)
	if err != nil {
		return err
	}

	tokenOut, _ := json.MarshalIndent(token, "", "  ")
	fmt.Println("token", string(tokenOut))

	claims := ClaimsFromToken(token)
	claimsOut, _ := json.MarshalIndent(claims, "", "  ")
	fmt.Println("claims", string(claimsOut))

	// FIXME: handle --format with text/template
	return nil
}
