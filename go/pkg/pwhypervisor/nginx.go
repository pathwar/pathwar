package pwhypervisor

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
	"pathwar.land/go/internal/randstring"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwcompose"
)

const (
	tmpNewNginxConfigFileName = "default.conf.new"
	nginxContainerName        = "pathwar-hypervisor-nginx"
)

func Nginx(ctx context.Context, opts HypervisorOpts, cli *client.Client, logger *zap.Logger) error {
	if opts.Salt == "" {
		opts.Salt = randstring.RandString(10)
		logger.Warn("random salt generated", zap.String("salt", opts.Salt))
	}
	if opts.ModeratorPassword == "" {
		opts.ModeratorPassword = randstring.RandString(10)
		logger.Warn("random moderator password generated", zap.String("password", opts.ModeratorPassword))
	}
	logger.Debug("hypervisor nginx", zap.Any("opts", opts))

	pwInfo, err := pwcompose.GetPathwarInfo(ctx, cli)
	if err != nil {
		return errcode.ErrHypervisorGetPathwarInfo.Wrap(err)
	}

	for _, flavor := range pwInfo.RunningFlavors {
		for _, instance := range flavor.Instances {
			for _, port := range instance.Ports {
				// FIXME: support non-standard ports using labels (later)
				if port.PublicPort != 0 {
					// add [hash(allowed_user+salt)].[domainsuffix] for each allowed user in config
				}
			}

			// add status-[instance].[domainsuffix] delivering a static file -> will be used to check if an instance is configured (i.e., for monitoring) in config

			// add moderator-[instance].[domainsuffix] with passphrase authentication -> used to check if an instance is working even without being in the AllowedUsers list in config
		}
	}

	// check if nginx server is started, build and start it if needed
	nginxContainerID, running, err := checkNginxContainer(ctx, cli)
	if err != nil {
		return errcode.ErrCheckNginxContainer.Wrap(err)
	}

	if opts.ForceRecreate && nginxContainerID != "" {
		logger.Debug("nginx container remove", zap.String("id", nginxContainerID))
		err := cli.ContainerRemove(ctx, nginxContainerID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			return errcode.ErrRemoveNginxContainer.Wrap(err)
		}
		nginxContainerID = ""
	}

	if nginxContainerID == "" {
		logger.Debug("build nginx container", zap.Any("opts", opts))
		nginxContainerID, err = buildNginxContainer(ctx, cli, opts)
		if err != nil {
			return errcode.ErrBuildNginxContainer.Wrap(err)
		}
		running = false
	}

	if !running {
		logger.Debug("start nginx container", zap.String("id", nginxContainerID))
		err = cli.ContainerStart(ctx, nginxContainerID, types.ContainerStartOptions{})
		if err != nil {
			return errcode.ErrStartNginxContainer.Wrap(err)
		}
	}

	// generate config
	configData := NginxConfigData{
		Upstreams: []NginxUpstream{
			{
				Name: "Sup",
				Host: "0.0.0.0",
				Port: "1337",
			},
			{
				Name: "Friend",
				Host: "0.0.0.0",
				Port: "1338",
				Hashes: []string{
					"123",
					"456",
				},
			},
		},
		Opts: opts,
	}
	configTemplate, err := template.New("nginx-config").Parse(NginxConfigTemplate)
	if err != nil {
		return errcode.ErrParsingTemplate.Wrap(err)
	}
	var buf bytes.Buffer
	configTemplate.Execute(&buf, configData)
	config := buf.Bytes()
	fmt.Println(string(config))

	buf.Reset()
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: filepath.Join("/etc/nginx/conf.d", tmpNewNginxConfigFileName),
		Mode: 0755,
		Size: int64(len(config)),
	})
	if err != nil {
		return errcode.ErrWriteConfigFileHeader.Wrap(err)
	}
	if _, err := tw.Write(config); err != nil {
		return errcode.ErrWriteConfigFile.Wrap(err)
	}
	err = tw.Close()
	if err != nil {
		return errcode.ErrCloseTarWriter.Wrap(err)
	}
	logger.Debug("copy nginx config into the container", zap.String("container-id", nginxContainerID))
	err = cli.CopyToContainer(ctx, nginxContainerID, "/", &buf, types.CopyToContainerOptions{})
	if err != nil {
		return errcode.ErrCopyNginxConfigToContainer.Wrap(err)
	}

	// check new nginx config
	args := []string{"nginx", "-c", filepath.Join("conf.d", tmpNewNginxConfigFileName)}
	logger.Debug("send nginx command", zap.Strings("args", args))
	res, err := nginxSendCommand(ctx, cli, nginxContainerID, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandNewConfigCheck.Wrap(err)
	}
	logger.Debug("nginx -c result", zap.String("res", string(res)))
	resultStr := string(res)
	if resultStr[len(resultStr)-11:] != "successful\n" {
		args := []string{"rm", filepath.Join("/etc/nginx/conf.d", tmpNewNginxConfigFileName)}
		logger.Debug("send nginx command", zap.Strings("args", args))
		_, err := nginxSendCommand(ctx, cli, nginxContainerID, args...)
		if err != nil {
			return errcode.ErrNginxSendCommandNewConfigRemove.Wrap(err)
		}
		return errcode.ErrNginxNewConfigCheckFailed.Wrap(err)
	}

	// replace nginx config with new one that we just generated
	args = []string{"mv", filepath.Join("/etc/nginx/conf.d", tmpNewNginxConfigFileName), "/etc/nginx/conf.d/default.conf"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	_, err = nginxSendCommand(ctx, cli, nginxContainerID, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandConfigReplace.Wrap(err)
	}

	// new config hot reload
	args = []string{"nginx", "-s", "reload"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	_, err = nginxSendCommand(ctx, cli, nginxContainerID, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandReloadConfig.Wrap(err)
	}

	// check config by sending it to NGINX container and call command for that

	// if config is OK, smart reload

	// bonus: generate a status page that contains information about the configuration

	// usage example:
	//   pathwar --debug hypervisor nginx '{"instance1": [123, 456], "instance2": [789]}'

	// 1: get pathwar info: pwcompose.GetPathwarInfo(ctx, cli)
	// 2: generate an entire nginx configuration file based on pathwar info and NginxConfig)
	//
	//    for each running instances we need to have:
	//      - a status-[instance].[domainsuffix] delivering a static file -> will be used to check if an instance is configured (i.e., for monitoring)
	//      - a moderator-[instance].[domainsuffix] with passphrase authentication -> used to check if an instance is working even without being in the AllowedUsers list
	//      - one [hash(allowed_user+salt)].[domainsuffix] without authentication for each allowed users of each instances

	// 3: check that new config file is valid -> there is an nginx command for that
	// 4: smart start/reload -> try to never shutdown the nginx server, if the server is already started, "nginx reload" can update the config, else we need to start it by ourselve
	// 5: bonus: generate a status page that contains information about the configuration

	return nil
}

// https://github.com/pathwar/pathwar/blob/v1.0.0/hypervisor/cmd_hypervisor_run.go
func buildNginxContainer(ctx context.Context, cli *client.Client, opts HypervisorOpts) (string, error) {
	out, err := cli.ImagePull(ctx, "docker.io/library/"+opts.NginxDockerImage, types.ImagePullOptions{})
	if err != nil {
		return "", errcode.ErrDockerAPIImagePull.Wrap(err)
	}
	io.Copy(os.Stdout, out)

	hostBinding := nat.PortBinding{
		HostIP:   opts.HostIP,
		HostPort: opts.HostPort,
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		return "", errcode.ErrNatPortOpening.Wrap(err)
	}
	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: opts.NginxDockerImage,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, nginxContainerName)
	if err != nil {
		return "", errcode.ErrDockerAPIContainerCreate.Wrap(err)
	}

	return cont.ID, nil
}

func startNginxContainer(ctx context.Context, cli *client.Client, containerID string) error {
	cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	fmt.Println("nginx container started: ", containerID)

	return nil
}

func checkNginxContainer(ctx context.Context, cli *client.Client) (string, bool, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return "", false, errcode.ErrDockerAPIContainerList.Wrap(err)
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if name[1:] == nginxContainerName {
				return container.ID, container.State == "running", nil
			}
		}
	}
	return "", false, nil
}

func nginxSendCommand(ctx context.Context, cli *client.Client, nginxContainerID string, args ...string) ([]byte, error) {
	cmd := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          args,
	}
	execID, err := cli.ContainerExecCreate(ctx, nginxContainerID, cmd)
	if err != nil {
		return nil, errcode.ErrDockerAPIContainerExecCreate.Wrap(err)
	}

	execConfig := types.ExecConfig{}
	res, err := cli.ContainerExecAttach(ctx, execID.ID, execConfig)
	if err != nil {
		return nil, errcode.ErrDockerAPIContainerExecAttach.Wrap(err)
	}

	err = cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, errcode.ErrDockerAPIContainerExecStart.Wrap(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Reader)
	return buf.Bytes(), nil
}
