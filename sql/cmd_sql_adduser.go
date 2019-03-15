package sql

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.pw/pkg/cli"
)

type adduserOptions struct {
	sql Options `mapstructure:"sql"`

	email    string `mapstructure:"email"`
	username string `mapstructure:"username"`
	password string `mapstructure:"password"`
}

type adduserCommand struct{ opts adduserOptions }

func (cmd *adduserCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "adduser",
		Args: func(_ *cobra.Command, args []string) error {
			if cmd.opts.email == "" {
				return errors.New("--email is mandatory")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			opts := cmd.opts
			opts.sql = GetOptions(commands)
			return runAdduser(opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	commands["sql"].ParseFlags(cc.Flags())
	return cc
}
func (cmd *adduserCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *adduserCommand) ParseFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&cmd.opts.email, "email", "", "", "valid email address")
	flags.StringVarP(&cmd.opts.username, "username", "", "", "random value if empty")
	flags.StringVarP(&cmd.opts.password, "password", "", "", "random value if empty")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runAdduser(opts adduserOptions) error {
	return fmt.Errorf("implementation is outdated and needs to be updated...")
	/*
		db, err := FromOpts(&opts.sql)
		if err != nil {
			return err
		}

		user := entity.User{
			Email:        opts.email,
			Username:     opts.username,
			PasswordSalt: "FIXME: randomize",
		}
		user.PasswordHash = "FIXME: generate"

		// FIXME: randomize username, password if empty
		// FIXME: verify email address validity
		// FIXME: verify email address spam/blacklist
		// FIXME: user.Validate()

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		out, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))

		return nil
	*/
}
