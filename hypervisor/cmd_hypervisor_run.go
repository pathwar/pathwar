package hypervisor

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.pw/pkg/cli"
)

type runOptions struct {
	target  string `mapstructure:"target"`
	webPort int    `mapstructure:"web-port"`
	// port, other options
	// driver=docker
}

type runCommand struct{ opts runOptions }

func (cmd *runCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "run",
		Args: func(_ *cobra.Command, args []string) error {
			if cmd.opts.target == "" {
				return errors.New("--target is mandatory")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			opts := cmd.opts
			return runRun(opts)
		},
	}
	cmd.ParseFlags(cc.Flags())
	return cc
}
func (cmd *runCommand) LoadDefaultOptions() error { return viper.Unmarshal(&cmd.opts) }
func (cmd *runCommand) ParseFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&cmd.opts.target, "target", "t", "", "target (image, path, etc)")
	flags.IntVarP(&cmd.opts.webPort, "web-port", "p", 8080, "web listening port")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runRun(opts runOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create docker client")
	}

	// configure new container
	imageInspect, _, err := cli.ImageInspectWithRaw(ctx, opts.target)
	if err != nil {
		return errors.Wrap(err, "failed to inspect image")
	}
	containerConfig := &container.Config{
		Image:     opts.target,
		Tty:       true,
		OpenStdin: true,
		// StdinOnce: true,
		// AttachStdin: true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: nat.PortSet{
			nat.Port("80/tcp"): {},
		},
		// Hostname: ""
		Cmd: imageInspect.Config.Cmd,
	}
	hostConfig := &container.HostConfig{
		// Binds: /etc/timezone
		PortBindings: nat.PortMap{
			nat.Port("80/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", opts.webPort)}},
		},
		AutoRemove: true,
		// RestartPolicy: "no"
		// Dns
		// DnsOptions
		// DnsSearch
	}
	//networkingConfig := &network.NetworkingConfig{SandboxID:"XXX",SandboxKey:"XXX"}

	// FIXME: create a limited network config
	// FIXME: restrict resources (cgroups, etc.)
	// FIXME: handle exposed port
	// FIXME: handle entrypoint/cmd
	// FIXME: wrap with custom init
	// FIXME: copy init binary into image
	// FIXME: configure env
	// FIXME: copy tokens somewhere safe in the image

	cont, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return err
	}
	if cont.Warnings != nil && len(cont.Warnings) > 0 {
		for _, warn := range cont.Warnings {
			zap.L().Warn(warn)
		}
	}

	if err := cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// FIXME: wait for error or ctrl+C
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, errC := cli.ContainerWait(ctxCancel, cont.ID, "")
	if err := <-errC; err != nil {
		return err
	}

	duration := 2 * time.Second
	if err := cli.ContainerStop(ctx, cont.ID, &duration); err != nil {
		return errors.Wrap(err, "failed to stop container")
	}

	if err := cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	}); err != nil {
		return errors.Wrap(err, "failed to remove container")
	}

	return nil
}
