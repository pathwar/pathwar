package sql

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func sqlSetupFlags(flags *pflag.FlagSet, opts *SQLOptions) {
	flags.StringVar(&opts.Path, "sql-path", "/tmp/pathwar.db", "SQL db path")
	viper.BindPFlags(flags)
}

func NewSQLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sql",
	}
	cmd.AddCommand(NewSQLDumpCommand())
	return cmd
}

func NewSQLDumpCommand() *cobra.Command {
	opts := &SQLOptions{}
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
