package main // import "pathwar.land"

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"pathwar.land/hypervisor"
	"pathwar.land/pkg/cli"
	"pathwar.land/server"
	"pathwar.land/sql"
	"pathwar.land/version"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rootCmd := newRootCommand()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		// Use: "pathwar.land",
		Use: os.Args[0],
	}
	cmd.PersistentFlags().BoolP("help", "h", false, "print usage")
	//cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	cmd.Version = fmt.Sprintf("%s (commit=%q, date=%q, built-by=%q)", version.Version, version.Commit, version.Date, version.BuiltBy)

	// Add commands
	commands := cli.Commands{}
	for name, command := range sql.Commands() {
		commands[name] = command
	}
	for name, command := range server.Commands() {
		commands[name] = command
	}
	for name, command := range hypervisor.Commands() {
		commands[name] = command
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// setup logging
		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(zap.DebugLevel)
		config.DisableStacktrace = true
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		l, err := config.Build()
		if err != nil {
			return errors.Wrap(err, "failed to configure logger")
		}
		zap.ReplaceGlobals(l)
		zap.L().Debug("logger initialized")

		// setup viper
		viper.AddConfigPath(".")
		viper.SetConfigName(".pathwar")
		viper.SetEnvPrefix("PATHWAR")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return errors.Wrap(err, "failed to apply viper config")
			}
		}

		for _, command := range commands {
			if err := command.LoadDefaultOptions(); err != nil {
				return err
			}
		}

		return nil
	}

	for name, command := range commands {
		if strings.Contains(name, " ") { // do not add commands where level > 1
			continue
		}
		cmd.AddCommand(command.CobraCommand(commands))
	}

	return cmd
}
