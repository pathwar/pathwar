package sql

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"pathwar.land/entity"
	"pathwar.land/pkg/cli"
	"pathwar.land/pkg/randstring"
)

type adduserOptions struct {
	sql Options `mapstructure:"sql"`

	email      string `mapstructure:"email"`
	username   string `mapstructure:"username"`
	locale     string `mapstructure:"locale"`
	password   string `mapstructure:"password"`
	websiteURL string `mapstructure:"website-url"`
	isStaff    bool   `mapstructure:"is-staff"`
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
	flags.StringVarP(&cmd.opts.locale, "locale", "", "fr_FR", "locale")
	flags.StringVarP(&cmd.opts.websiteURL, "website-url", "", "", "website url")
	flags.BoolVarP(&cmd.opts.isStaff, "is-staff", "", false, "is staff?")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runAdduser(opts adduserOptions) error {
	db, err := FromOpts(&opts.sql)
	if err != nil {
		return err
	}

	if opts.password == "" {
		opts.password = randstring.RandString(15)
		zap.L().Info("password is empty, generating a new one", zap.String("password", opts.password))
	}
	if opts.username == "" {
		opts.username = randstring.RandString(10)
		zap.L().Info("username is empty, generating a new one", zap.String("username", opts.username))
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(opts.password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := entity.User{
		Username:   opts.username,
		WebsiteURL: opts.websiteURL,
		IsStaff:    opts.isStaff,
		Locale:     opts.locale,
		AuthMethods: []*entity.AuthMethod{
			{
				Identifier:   opts.email,
				EmailAddress: opts.email,
				PasswordHash: string(hash),
				Provider:     entity.AuthMethod_EmailAndPassword,
				IsVerified:   true,
			},
		},
	}

	// FIXME: verify email address validity
	// FIXME: verify email address spam/blacklist
	// FIXME: verify email for duplicate
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
}
