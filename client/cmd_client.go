package client

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.land/pkg/cli"
)

type Options struct {
	Token     string `mapstructure:"token"`
	Client    string `mapstructure:"client"`
	Realm     string `mapstructure:"realm"`
	AuthURL   string `mapstructure:"auth-url"`
	PublicKey string `mapstructure:"public-key"`
}

func Commands() cli.Commands {
	return cli.Commands{
		"client":        &clientCommand{},
		"client whoami": &whoamiCommand{},
		"client logout": &logoutCommand{},
		"client claims": &claimsCommand{},
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
	flags.StringVarP(&cmd.opts.PublicKey, "public-key", "", "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlEFxLlywsbI5BQ7DVkA66fICWGIYPpD+aZNYRR7SIc0zdtJR4xMOt5CjM0vbYT4z2a1U2yl0ewunyxFm8niS8w6mKYFnOS4nnSchQyIAmJkpLC4eAjijCdEHdr8mSqamThSrVRGSYEEsa+adidC13kRDy7NDKhvZb8F0YqnktNk6WHSlb8r2QRLPJ1DX534jjXPY6l/eoHuLJAOZxBlfwV5Dg37TVmf2xAH812E7ZigycLAvhsMvr5x2jLavAEEnZZmlQf4cyQ4tlMzKS1Zp0NcdOGS/i6lrndc5pNtZQuGr8IGBrEbTRFUiavn/HDnyalYZy8T5LakXRdVaKdshAQIDAQAB", "SSO Public Key")
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
	command.AddCommand(commands["client claims"].CobraCommand(commands))
	return command
}
