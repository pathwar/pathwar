package hypervisor

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.land/pkg/cli"
)

type Options struct {
}

func (opts Options) String() string {
	out, _ := json.Marshal(opts)
	return string(out)
}

func Commands() cli.Commands {
	return cli.Commands{
		"hypervisor":       &hypervisorCommand{},
		"hypervisor run":   &runCommand{},
		"hypervisor prune": &pruneCommand{},
	}
}

type hypervisorCommand struct{ opts Options }

func (cmd *hypervisorCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *hypervisorCommand) ParseFlags(flags *pflag.FlagSet) {
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}
func (cmd *hypervisorCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	command := &cobra.Command{
		Use: "hypervisor",
	}
	command.AddCommand(commands["hypervisor run"].CobraCommand(commands))
	command.AddCommand(commands["hypervisor prune"].CobraCommand(commands))
	return command
}
