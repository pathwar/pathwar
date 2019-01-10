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

func NewSQLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sql",
	}
	cmd.AddCommand(NewSQLDumpCommand())
	return cmd
}

func NewSQLDumpCommand() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}
			return sqlDump(opts)
		},
	}
	sqlSetupFlags(cmd.Flags(), opts)
	return cmd
}
