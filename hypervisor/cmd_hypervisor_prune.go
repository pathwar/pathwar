package hypervisor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"pathwar.land/pkg/cli"
)

type pruneOptions struct {
	// timeout
	// driver=docker
}

type pruneCommand struct{ opts pruneOptions }

func (cmd *pruneCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "prune",
		RunE: func(_ *cobra.Command, args []string) error {
			opts := cmd.opts
			return runPrune(opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	return cc
}
func (cmd *pruneCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *pruneCommand) ParseFlags(flags *pflag.FlagSet) {
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runPrune(opts pruneOptions) error {
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
	log.Printf("%d container(s) stopped", len(containers))

	// stop containers
	var g errgroup.Group
	var timeout = 5 * time.Second
	for _, container := range containers {
		g.Go(func() error {
			return cli.ContainerStop(ctx, container.ID, &timeout)
		})
	}
	if err := g.Wait(); err != nil {
		log.Printf("some container failed to stop: %v", err)
	}

	// prune containers
	report, err := cli.ContainersPrune(ctx, filters)
	if err != nil {
		return err
	}
	if len(report.ContainersDeleted) > 0 {
		log.Printf("%d container(s) pruned", len(report.ContainersDeleted))
	}
	return err
}
