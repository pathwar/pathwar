package hypervisor

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.land/pkg/cli"
)

type psOptions struct {
	NoTrunc bool
	// timeout
	// driver=docker
}

type psCommand struct{ opts psOptions }

func (cmd *psCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "ps",
		RunE: func(_ *cobra.Command, args []string) error {
			opts := cmd.opts
			return runPs(opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	return cc
}
func (cmd *psCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *psCommand) ParseFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&cmd.opts.NoTrunc, "no-trunc", "", false, "no truncate")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runPs(opts psOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create docker client")
	}

	filters := filters.NewArgs()
	filters.Add("label", fmt.Sprintf("%s", createdByPathwarLabel))

	// list containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters,
	})
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	// FIXME: support non-docker instances
	defer w.Flush()
	fmt.Fprintln(w, "ID\tIMAGE\tCREATED\tSTATUS\tPORTS\tNAME")
	for _, container := range containers {
		shortID := truncIf(container.ID, 8, !opts.NoTrunc)
		createdAt := units.HumanDuration(time.Now().UTC().Sub(time.Unix(container.Created, 0)))

		defer w.Flush()
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			shortID,
			container.Image,
			createdAt,
			container.State,
			"not implemmented", // need to call ContainerExecInspect on running containers to get the information
			container.Names[0],
		)
	}

	return err
}

func truncIf(input string, length int, trunc bool) string {
	if !trunc || len(input) < length {
		return input
	}
	return input[:length]
}
