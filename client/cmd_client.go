package client

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.land/pkg/cli"
)

type Options struct {
	Token   string `mapstructure:"token"`
	Client  string `mapstructure:"client"`
	Realm   string `mapstructure:"realm"`
	AuthURL string `mapstructure:"auth-url"`
}

func Commands() cli.Commands {
	return cli.Commands{
		"client":        &clientCommand{},
		"client whoami": &whoamiCommand{},
		"client logout": &logoutCommand{},
	}
}

func GetOptions(commands cli.Commands) Options {
	return commands["client"].(*clientCommand).opts
}

type clientCommand struct{ opts Options }

func (cmd *clientCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *clientCommand) ParseFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&cmd.opts.Token, "token", "", "", "SSO Token")
	flags.StringVarP(&cmd.opts.Client, "sso-client", "", "platform-cli", "SSO Client")
	flags.StringVarP(&cmd.opts.Realm, "realm", "", "Pathwar-Dev", "SSO Realm")
	flags.StringVarP(&cmd.opts.AuthURL, "auth-url", "", "https://sso.pathwar.land/auth", "SSO Authentication URL")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}
func (cmd *clientCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	command := &cobra.Command{
		Use: "client",
	}
	command.AddCommand(commands["client whoami"].CobraCommand(commands))
	command.AddCommand(commands["client logout"].CobraCommand(commands))
	return command
}
