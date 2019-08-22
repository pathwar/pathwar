package sql

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"pathwar.land/pkg/cli"
)

type Options struct {
	Config string `mapstructure:"config"`
}

func Commands() cli.Commands {
	return cli.Commands{
		"sql":         &sqlCommand{},
		"sql dump":    &dumpCommand{},
		"sql adduser": &adduserCommand{},
		"sql info":    &infoCommand{},
	}
}

func GetOptions(commands cli.Commands) Options {
	return commands["sql"].(*sqlCommand).opts
}

type sqlCommand struct{ opts Options }

func (cmd *sqlCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *sqlCommand) ParseFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&cmd.opts.Config, "sql-config", "", "root:uns3cur3@tcp(127.0.0.1:3306)/pathwar?charset=utf8&parseTime=true", "SQL connection config")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}
func (cmd *sqlCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	command := &cobra.Command{
		Use: "sql",
	}
	command.AddCommand(commands["sql dump"].CobraCommand(commands))
	command.AddCommand(commands["sql info"].CobraCommand(commands))
	command.AddCommand(commands["sql adduser"].CobraCommand(commands))
	return command
}
