package pwagent

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/martinlindhe/base36"
	"github.com/moby/moby/pkg/stdcopy"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"pathwar.land/go/internal/randstring"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwcompose"
)

const nginxConfigTemplate = `
{{$root := .}}
#user                 www www;
worker_processes     5;
error_log            /proc/self/fd/2;
#pid                  /tmp/nginx.pid;
worker_rlimit_nofile 8192;

events {
  worker_connections 4096;
}

http {
  types {
    text/html                             html htm shtml;
    text/css                              css;
    image/gif                             gif;
    image/jpeg                            jpeg jpg;
    application/x-javascript              js;
    text/plain                            txt;
    image/png                             png;
    image/x-icon                          ico;
  }
  proxy_redirect          off;
  proxy_set_header        Host            $host;
  proxy_set_header        X-Real-IP       $remote_addr;
  proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
  client_max_body_size    10m;
  client_body_buffer_size 128k;
  proxy_connect_timeout   90;
  proxy_send_timeout      90;
  proxy_read_timeout      90;
  proxy_buffers           32 4k;
  index                   index.html index.htm;

  default_type                  application/octet-stream;
  log_format                    main '$remote_addr - $remote_user [$time_local]  $status "$request" $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"';
  access_log                    /proc/self/fd/1 main;
  sendfile                      on;
  tcp_nopush                    on;
  server_names_hash_bucket_size 128;

  server {
    listen      80 default_server;
    server_name _;
    error_log   /proc/self/fd/2;
    access_log  /proc/self/fd/1;
    return      503;
  }

  {{range .Upstreams -}}
  upstream upstream_{{.Name}} { server {{.Host}}:{{.Port}}; }
  server {
    listen      80;
    server_name moderator-{{.Name}}{{$root.Opts.DomainSuffix}};
    access_log  /proc/self/fd/1;
    error_log   /proc/self/fd/2;
    # FIXME: add auth
    location / {
      proxy_pass http://upstream_{{.Name}};
    }
  }
  {{- if not (eq (len .Hashes) 0) }}
  server {
    listen      80;
    server_name{{range .Hashes}} {{.}}{{$root.Opts.DomainSuffix}}{{end}};
    access_log  /proc/self/fd/1;
    error_log   /proc/self/fd/2;
    location / {
      proxy_pass http://upstream_{{.Name}};
    }
  }
  {{end}}
  {{end -}}
}
`

func Nginx(ctx context.Context, opts AgentOpts, cli *client.Client, logger *zap.Logger) error {
	if opts.Salt == "" {
		opts.Salt = randstring.RandString(10)
		logger.Warn("random salt generated", zap.String("salt", opts.Salt))
	}
	if opts.ModeratorPassword == "" {
		opts.ModeratorPassword = randstring.RandString(10)
		logger.Warn("random moderator password generated", zap.String("password", opts.ModeratorPassword))
	}
	logger.Debug("agent nginx", zap.Any("opts", opts))

	// check if proxy network has been created
	proxyNetworkID := ""
	networkResources, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return errcode.ErrDockerAPINetworkList.Wrap(err)
	}
	for _, networkResource := range networkResources {
		if networkResource.Name == pwcompose.ProxyNetworkName {
			proxyNetworkID = networkResource.ID
		}
	}
	if proxyNetworkID == "" {
		logger.Debug("proxy network create", zap.String("name", pwcompose.ProxyNetworkName))
		response, err := cli.NetworkCreate(ctx, pwcompose.ProxyNetworkName, types.NetworkCreate{
			CheckDuplicate: true,
		})
		proxyNetworkID = response.ID
		if err != nil {
			return errcode.ErrDockerAPINetworkCreate.Wrap(err)
		}
	}

	// check if nginx server is started
	nginxContainer, err := checkNginxContainer(ctx, cli)
	if err != nil {
		return errcode.ErrCheckNginxContainer.Wrap(err)
	}

	// remove nginx container if forced
	if opts.ForceRecreate && nginxContainer != nil {
		logger.Debug("nginx container remove", zap.String("id", nginxContainer.ID))
		err := cli.ContainerRemove(ctx, nginxContainer.ID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			return errcode.ErrRemoveNginxContainer.Wrap(err)
		}
		nginxContainer = nil
	}

	// build nginx container if needed
	var nginxContainerID string
	running := false
	if nginxContainer == nil {
		logger.Debug("build nginx container", zap.Any("opts", opts))
		nginxContainerID, err = buildNginxContainer(ctx, cli, opts)
		if err != nil {
			return errcode.ErrBuildNginxContainer.Wrap(err)
		}
	} else {
		nginxContainerID = nginxContainer.ID
		running = (nginxContainer.State == "running")
	}

	// start nginx container if needed
	if !running {
		logger.Debug("start nginx container", zap.String("id", nginxContainerID))
		err = cli.ContainerStart(ctx, nginxContainerID, types.ContainerStartOptions{})
		if err != nil {
			return errcode.ErrStartNginxContainer.Wrap(err)
		}
		nginxContainer, err = checkNginxContainer(ctx, cli)
		if err != nil {
			return errcode.ErrCheckNginxContainer.Wrap(err)
		}
	}

	// connect nginx container to proxy network
	if _, onProxyNetwork := nginxContainer.NetworkSettings.Networks[pwcompose.ProxyNetworkName]; !onProxyNetwork {
		logger.Debug("connect nginx to proxy network", zap.String("nginx-id", nginxContainer.ID), zap.String("network-id", proxyNetworkID))
		err = cli.NetworkConnect(ctx, proxyNetworkID, nginxContainer.ID, nil)
		if err != nil {
			return errcode.ErrNginxConnectNetwork.Wrap(err)
		}
		// refresh container struct so it contains network configuration
		nginxContainer, err = checkNginxContainer(ctx, cli)
		if err != nil {
			return errcode.ErrCheckNginxContainer.Wrap(err)
		}
	}

	// update proxy network nginx container IP
	proxyNetworkIP := nginxContainer.NetworkSettings.Networks[pwcompose.ProxyNetworkName].IPAddress

	// make sure that exposed containers are connected to proxy network
	containersInfo, err := pwcompose.GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrAgentGetContainersInfo.Wrap(err)
	}
	for _, flavor := range containersInfo.RunningFlavors {
		for _, instance := range flavor.Instances {
			for _, port := range instance.Ports {
				if port.PrivatePort != 0 {
					if _, onProxyNetwork := instance.NetworkSettings.Networks[pwcompose.ProxyNetworkName]; !onProxyNetwork {
						logger.Debug("connect container to proxy network", zap.String("container-id", instance.ID), zap.String("network-id", proxyNetworkID))
						err = cli.NetworkConnect(ctx, proxyNetworkID, instance.ID, nil)
						if err != nil {
							return errcode.ErrContainerConnectNetwork.Wrap(err)
						}
					}
					break
				}
			}
		}
	}

	// update domainsuffix for local use if needed
	if opts.DomainSuffix == "local" {
		opts.DomainSuffix = "." + proxyNetworkIP + ".xip.io"
	}

	configData := NginxConfigData{
		Upstreams: []NginxUpstream{},
		Opts:      opts,
	}

	// update config data with containers infos
	containersInfo, err = pwcompose.GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrAgentGetContainersInfo.Wrap(err)
	}
	for _, flavor := range containersInfo.RunningFlavors {
		for _, instance := range flavor.Instances {
			for _, port := range instance.Ports {
				// FIXME: support non-standard ports using labels (later)
				upstream := NginxUpstream{
					Hashes: []string{},
				}
				if port.PrivatePort != 0 {
					upstream.Name = instance.Names[0][1:]
					upstream.Host = instance.NetworkSettings.Networks[pwcompose.ProxyNetworkName].IPAddress
					upstream.Port = strconv.Itoa(int(port.PrivatePort))

					// add hash per users to proxy configuration
					if _, found := opts.AllowedUsers[instance.Names[0][1:]]; found {
						for _, userID := range opts.AllowedUsers[instance.Names[0][1:]] {
							hash, err := generatePrefixHash(instance.ID, userID, opts.Salt)
							if err != nil {
								return errcode.ErrGeneratePrefixHash.Wrap(err)
							}
							upstream.Hashes = append(upstream.Hashes, hash)
						}
					}
					configData.Upstreams = append(configData.Upstreams, upstream)
					// FIXME: doesn't handle multiple port per instance yet
					break
				}
			}
		}
	}

	buf, err := buildNginxConfigTar(configData)
	if err != nil {
		return errcode.ErrBuildNginxConfig.Wrap(err)
	}

	logger.Debug("copy nginx config into the container", zap.String("container-id", nginxContainer.ID))
	err = cli.CopyToContainer(ctx, nginxContainer.ID, "/etc/nginx/", buf, types.CopyToContainerOptions{})
	if err != nil {
		return errcode.ErrCopyNginxConfigToContainer.Wrap(err)
	}

	// check new nginx config
	args := []string{"nginx", "-t", "-c", "/etc/nginx/nginx.conf"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	err = nginxSendCommand(ctx, cli, nginxContainer.ID, logger, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandNewConfigCheck.Wrap(err)
	}

	// new config hot reload
	args = []string{"nginx", "-s", "reload"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	err = nginxSendCommand(ctx, cli, nginxContainer.ID, logger, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandReloadConfig.Wrap(err)
	}

	for _, upstream := range configData.Upstreams {
		for _, hash := range upstream.Hashes {
			fmt.Println(upstream.Name + ": " + hash + opts.DomainSuffix)
		}

	}

	return nil
}

func buildNginxConfigTar(data interface{}) (*bytes.Buffer, error) {
	configTemplate, err := template.New("nginx-config").Parse(nginxConfigTemplate)
	if err != nil {
		return nil, errcode.ErrParsingTemplate.Wrap(err)
	}
	var configBuf bytes.Buffer
	err = configTemplate.Execute(&configBuf, data)
	if err != nil {
		return nil, errcode.ErrExecuteTemplate.Wrap(err)
	}
	config := configBuf.Bytes()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: "nginx.conf",
		Mode: 0755,
		Size: int64(len(config)),
	})
	if err != nil {
		return nil, errcode.ErrWriteConfigFileHeader.Wrap(err)
	}

	if _, err := tw.Write(config); err != nil {
		return nil, errcode.ErrWriteConfigFile.Wrap(err)
	}

	err = tw.Close()
	if err != nil {
		return nil, errcode.ErrCloseTarWriter.Wrap(err)
	}

	return &buf, nil
}

// https://github.com/pathwar/pathwar/blob/v1.0.0/agent/cmd_agent_run.go
func buildNginxContainer(ctx context.Context, cli *client.Client, opts AgentOpts) (string, error) {
	out, err := cli.ImagePull(ctx, opts.NginxDockerImage, types.ImagePullOptions{})
	if err != nil {
		return "", errcode.ErrDockerAPIImagePull.Wrap(err)
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return "", errcode.TODO.Wrap(err)
	}

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
		}, nil, pwcompose.NginxContainerName)
	if err != nil {
		return "", errcode.ErrDockerAPIContainerCreate.Wrap(err)
	}

	return cont.ID, nil
}

func checkNginxContainer(ctx context.Context, cli *client.Client) (*types.Container, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, errcode.ErrDockerAPIContainerList.Wrap(err)
	}
	for _, container := range containers {
		for _, name := range container.Names {
			if name[1:] == pwcompose.NginxContainerName {
				return &container, nil
			}
		}
	}
	return nil, nil
}

func nginxSendCommand(ctx context.Context, cli *client.Client, nginxContainerID string, logger *zap.Logger, args ...string) error {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          args,
	}
	execRes, err := cli.ContainerExecCreate(ctx, nginxContainerID, execConfig)
	if err != nil {
		return errcode.ErrDockerAPIContainerExecCreate.Wrap(err)
	}

	res, err := cli.ContainerExecAttach(ctx, execRes.ID, execConfig)
	if err != nil {
		return errcode.ErrDockerAPIContainerExecAttach.Wrap(err)
	}
	defer res.Close()

	var outbuf, errbuf bytes.Buffer
	outputDone := make(chan error)
	go func() {
		_, err = stdcopy.StdCopy(&outbuf, &errbuf, res.Reader)
		outputDone <- err
	}()
	select {
	case err := <-outputDone:
		if err != nil {
			return errcode.TODO.Wrap(err)
		}
		break
	case <-ctx.Done():
		return errcode.TODO.Wrap(ctx.Err())
	}
	stdout, err := ioutil.ReadAll(&outbuf)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	stderr, err := ioutil.ReadAll(&errbuf)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	inspect, err := cli.ContainerExecInspect(ctx, execRes.ID)
	if err != nil {
		return errcode.ErrDockerAPIContainerExecInspect.Wrap(err)
	}

	if stderr != nil && len(stderr) > 0 {
		logger.Warn("exec finished with stderr", zap.String("stderr", string(stderr)))
	}
	logger.Debug("exec finished",
		zap.String("stdout", string(stdout)),
		zap.Int("exit-code", inspect.ExitCode),
		zap.Bool("running", inspect.Running),
	)

	if inspect.ExitCode != 0 {
		return errcode.ErrDockerAPIExitCode
	}

	return nil
}

func generatePrefixHash(instanceID string, userID int64, salt string) (string, error) {
	stringToHash := fmt.Sprintf("%s%d%s", instanceID, userID, salt)
	hashBytes := make([]byte, 8)
	hasher := sha3.NewShake256()
	_, err := hasher.Write([]byte(stringToHash))
	if err != nil {
		return "", errcode.ErrWriteBytesToHashBuilder.Wrap(err)
	}
	_, err = hasher.Read(hashBytes)
	if err != nil {
		return "", errcode.ErrReadBytesFromHashBuilder.Wrap(err)
	}
	userHash := strings.ToLower(base36.EncodeBytes(hashBytes))[:8] // we voluntarily expect short hashes here
	return userHash, nil
}
