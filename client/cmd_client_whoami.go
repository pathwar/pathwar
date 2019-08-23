package client

import (
	"errors"
	"fmt"

	"github.com/keycloak/kcinit/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.land/pkg/cli"
)

type whoamiOptions struct{ client Options }
type whoamiCommand struct{ opts whoamiOptions }

func (cmd *whoamiCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "whoami",
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
			return runWhoami(&cmd.opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	commands["client"].ParseFlags(cc.Flags())
	return cc
}
func (cmd *whoamiCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *whoamiCommand) ParseFlags(flags *pflag.FlagSet) {
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runWhoami(opts *whoamiOptions) error {
	keycloak := rest.New()
	realmURL := fmt.Sprintf("%s/realms/%s", opts.client.AuthURL, opts.client.Realm)
	base := keycloak.Target(realmURL)
	if base == nil {
		return errors.New("failed to initialize keycloak client")
	}
	oidc := base.Path("protocol/openid-connect")
	res, err := oidc.Path("userinfo").Request().Header("Authorization", "bearer "+opts.client.Token).Get()
	if err != nil {
		return err
	}
	var info map[string]interface{}
	if err := res.ReadJson(&info); err != nil {
		return err
	}
	for k, v := range info {
		fmt.Printf("- %s: %v\n", k, v)
	}
	return nil
}
