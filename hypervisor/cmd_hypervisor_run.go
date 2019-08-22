package hypervisor

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"pathwar.land/pkg/cli"
	"pathwar.land/pkg/randstring"
	pwctlconfig "pathwar.land/pwctl/config"
)

type runOptions struct {
	detach            bool   `mapstructure:"detach"`
	nginxProxy        bool   `mapstructure:"nginx-proxy"`
	nginxProxyHost    string `mapstructure:"nginx-proxy-host"`
	nginxProxyNetwork string `mapstructure:"nginx-proxy-network"`
	override          bool   `mapstructure:"override"`
	target            string `mapstructure:"target"`
	webPort           int    `mapstructure:"web-port"`
	// tcpPort, udpPort
	// pull
	// driver=docker
}

const createdByPathwarLabel = "created-by-pathwar-hypervisor"

type runCommand struct{ opts runOptions }

func (cmd *runCommand) CobraCommand(commands cli.Commands) *cobra.Command {
	cc := &cobra.Command{
		Use: "run",
		Args: func(_ *cobra.Command, args []string) error {
			switch {
			case len(args) == 0 && cmd.opts.target == "":
				return errors.New("--target is mandatory")
			case len(args) == 1 && cmd.opts.target == "":
				cmd.opts.target = args[0]
			case len(args) == 0 && cmd.opts.target != "":
				// everything is okay!
			default:
				return errors.New("bad usage")
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
	flags.BoolVarP(&cmd.opts.detach, "detach", "d", false, "detach mode")
	flags.BoolVarP(&cmd.opts.nginxProxy, "nginx-proxy", "", false, "use nginx-proxy instead of exposing web port")
	flags.BoolVarP(&cmd.opts.override, "override", "", false, "prune existing repo if existing")
	flags.IntVarP(&cmd.opts.webPort, "web-port", "p", -1, "web listening port (random if unset)")
	flags.StringVarP(&cmd.opts.nginxProxyHost, "nginx-proxy-host", "", "", "host name for nginx-proxy (use nip.io if empty)")
	flags.StringVarP(&cmd.opts.nginxProxyNetwork, "nginx-proxy-network", "", "service-proxy", "network name for nginx-proxy")
	flags.StringVarP(&cmd.opts.target, "target", "t", "", "target (image, path, etc)")
	if err := viper.BindPFlags(flags); err != nil {
		zap.L().Warn("failed to bind viper flags", zap.Error(err))
	}
}

func runRun(opts runOptions) error {
	// FIXME: ensure proxy is setup with nip.io as default hostname
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create docker client")
	}
	zap.L().Debug("connected to Docker daemon",
		zap.String("host", cli.DaemonHost()),
		zap.String("version", cli.ClientVersion()),
	)

	// configure new container
	zap.L().Debug("inspecting image", zap.String("target", opts.target))
	imageInspect, _, err := cli.ImageInspectWithRaw(ctx, opts.target)
	if err != nil {
		return errors.Wrap(err, "failed to inspect image")
	}
	env := []string{"PATHWAR_HYPERVISOR=1"} // FIXME: use version
	if opts.nginxProxy {
		if opts.nginxProxyHost == "" {
			consensus := externalip.DefaultConsensus(nil, nil)
			ip, err := consensus.ExternalIP()
			if err != nil {
				return errors.Wrap(err, "failed to guess external IP, please specify --nginx-proxy-host instead")
			}
			opts.nginxProxyHost = fmt.Sprintf(
				"%s.%s.nip.io",
				strings.ToLower(randstring.RandString(10)),
				ip.String(),
			)
			zap.L().Info("container is configured to use nginx-proxy",
				zap.String("host", fmt.Sprintf("http://%s", opts.nginxProxyHost)),
			)
		}
		env = append(env, fmt.Sprintf("VIRTUAL_HOST=%s", opts.nginxProxyHost))
	}
	containerWebPort := nat.Port("80/tcp")
	exposedPorts := nat.PortSet{}
	exposedPorts[containerWebPort] = struct{}{} // FIXME: autodetect source port or allow override
	containerConfig := &container.Config{
		Image:        opts.target,
		Tty:          true,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: exposedPorts,
		Env:          env,
		Entrypoint:   strslice.StrSlice{"/bin/pwctl", "entrypoint"},
		Cmd:          append(imageInspect.Config.Entrypoint, imageInspect.Config.Cmd...),
		Labels:       map[string]string{createdByPathwarLabel: "true"},
		// StdinOnce: true,
		// AttachStdin: true,
		// Hostname: ""
	}
	hostPort := fmt.Sprintf("%d", opts.webPort)
	if opts.webPort == -1 {
		hostPort = ""
	}
	portBindings := nat.PortMap{}
	if !opts.nginxProxy {
		portBindings[containerWebPort] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: hostPort},
		}
	}
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        []string{"/etc/timezone:/etc/timezone:ro"},
		AutoRemove:   true,
		// Dns
		// DnsOptions
		// DnsSearch
	}
	if opts.detach {
		hostConfig.RestartPolicy = container.RestartPolicy{Name: "unless-stopped", MaximumRetryCount: 42}
		hostConfig.AutoRemove = false
	}
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	if opts.nginxProxy {
		networkingConfig.EndpointsConfig[opts.nginxProxyNetwork] = &network.EndpointSettings{}
	}
	// FIXME: create a limited network config
	// FIXME: restrict resources (cgroups, etc.)

	zap.L().Debug("creating container",
		zap.Any("container-config", containerConfig),
		zap.Any("host-config", hostConfig),
	)
	cont, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, "")
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
			randstring.RandString(10),
			randstring.RandString(10),
			randstring.RandString(10),
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
	zap.L().Debug("injecting pwctl into the container", zap.String("container-id", cont.ID))
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
	zap.L().Debug("starting container", zap.String("container-id", cont.ID))
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
	zap.L().Debug("connecting to container logs", zap.String("container-id", cont.ID))
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
	zap.L().Debug("stopping container", zap.String("container-id", cont.ID), zap.Duration("timeout", duration))
	if err := cli.ContainerStop(ctx, cont.ID, &duration); err != nil {
		return errors.Wrap(err, "failed to stop container")
	}

	zap.L().Debug("removing container", zap.String("container-id", cont.ID))
	if err := cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	}); err != nil {
		return errors.Wrap(err, "failed to remove container")
	}

	return nil
}
