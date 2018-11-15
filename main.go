package main // import "pathwar.pw"

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

	"pathwar.pw/server"
	"pathwar.pw/sql"
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
		Use: "pathwar.pw",
	}
	cmd.PersistentFlags().BoolP("help", "h", false, "print usage")
	//cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

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
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return errors.Wrap(err, "failed to apply viper config")
			}
		}

		return nil
	}

	cmd.AddCommand(
		server.NewServerCommand(),
		sql.NewSQLCommand(),
	)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}
