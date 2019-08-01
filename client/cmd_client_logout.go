package client

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/keycloak/kcinit/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.land/pkg/cli"
)

type logoutOptions struct{ client Options }
type logoutCommand struct{ opts logoutOptions }

func (cmd *logoutCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "logout",
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
			return runLogout(&cmd.opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	commands["client"].ParseFlags(cc.Flags())
	return cc
}
func (cmd *logoutCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *logoutCommand) ParseFlags(flags *pflag.FlagSet) {
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runLogout(opts *logoutOptions) error {
	keycloak := rest.New()
	realmURL := fmt.Sprintf("%s/realms/%s", opts.client.AuthURL, opts.client.Realm)
	base := keycloak.Target(realmURL)
	if base == nil {
		return errors.New("failed to initialize keycloak client")
	}
	oidc := base.Path("protocol/openid-connect")
	form := url.Values{}
	form.Set("client_id", opts.client.Client)
	// form.Set("client_secret", opts.client.Secret)
	form.Set("refresh_token", opts.client.Token)
	res, err := oidc.Path("logout").Request().Form(form).Post()
	if err != nil {
		return err
	}
	txt, err := res.ReadText()
	if err != nil {
		return err
	}
	fmt.Println(txt)
	return nil
}
