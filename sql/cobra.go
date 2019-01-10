package sql

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func sqlSetupFlags(flags *pflag.FlagSet, opts *Options) {
	flags.StringVar(&opts.Path, "sql-path", "/tmp/pathwar.db", "SQL db path")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

var globalOpts Options

func GetOptions() *Options {
	opts := globalOpts
	return &opts
}

func NewSQLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sql",
	}
	cmd.AddCommand(NewSQLDumpCommand())
	return cmd
}

func NewSQLDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sqlDump(GetOptions())
		},
	}
	sqlSetupFlags(cmd.Flags(), &globalOpts)
	return cmd
}
