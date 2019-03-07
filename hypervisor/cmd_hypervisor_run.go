package hypervisor

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.pw/pkg/cli"
	pwctlconfig "pathwar.pw/pwctl/config"
)

type runOptions struct {
	target  string `mapstructure:"target"`
	webPort int    `mapstructure:"web-port"`
	detach  bool   `mapstructure:"detach"`
	// port, other options
	// driver=docker
}

const createdByPathwarLabel = "created-by-pathwar-hypervisor"

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
	flags.BoolVarP(&cmd.opts.detach, "detach", "d", false, "detach mode")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runRun(opts runOptions) error {
	// FIXME: ensure proxy is setup with xip.io as default hostname
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
		Entrypoint: strslice.StrSlice{"/bin/pwctl", "entrypoint"},
		Cmd:        append(imageInspect.Config.Entrypoint, imageInspect.Config.Cmd...),
		Labels:     map[string]string{createdByPathwarLabel: "true"},
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
	// FIXME: configure env

	cont, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return err
	}
	if cont.Warnings != nil && len(cont.Warnings) > 0 {
		for _, warn := range cont.Warnings {
			zap.L().Warn(warn)
		}
	}

	pwctlConfig := pwctlconfig.Config{
		Passphrases: []string{
			randString(10),
			randString(10),
			randString(10),
		},
	}
	// if !pwctlConfig.Validate() ...
	pwctlConfigJSON, _ := json.Marshal(pwctlConfig)

	// inject tools & config in container
	// FIXME: support alternative architectures -> using copyFromContainer with a dedicated image?
	var buf bytes.Buffer
	var pwctlBox = packr.New("pwctl-binaries", "../pwctl/out")
	binary, err := pwctlBox.Find("pwctl-linux-amd64")
	if err != nil {
		return err
	}
	tw := tar.NewWriter(&buf)
	if err := tw.WriteHeader(&tar.Header{
		Name: "/bin/pwctl",
		Mode: 0755,
		Size: int64(len(binary)),
	}); err != nil {
		return err
	}
	if _, err := tw.Write(binary); err != nil {
		return err
	}
	if err := tw.WriteHeader(&tar.Header{
		Name: "/pwctl.json",
		Mode: 0755,
		Size: int64(len(pwctlConfigJSON)),
		// FIXME: chown it to container's default user
	}); err != nil {
		return err
	}
	if _, err := tw.Write(pwctlConfigJSON); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}
	if err := cli.CopyToContainer(
		ctx,
		cont.ID,
		"/",
		&buf,
		types.CopyToContainerOptions{},
	); err != nil {
		return err
	}

	// start container
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := cli.ContainerStart(ctxCancel, cont.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if opts.detach {
		fmt.Println(cont.ID)
		return nil
	}

	// FIXME: wait for error or ctrl+C
	/*_, errC := cli.ContainerWait(ctx, cont.ID, "")
	if err := <-errC; err != nil {
		return err
	}*/

	// read logs
	stream, err := cli.ContainerLogs(ctx, cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}
	if _, err = io.Copy(os.Stderr, stream); err != nil && err != io.EOF {
		return err
	}

	// cleanup
	duration := 50 * time.Millisecond
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
