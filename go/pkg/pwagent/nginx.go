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
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/pkg/stdcopy"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwcompose"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func applyNginxConfig(ctx context.Context, apiInstances *pwapi.AgentListInstances_Output, dockerClient *client.Client, opts Opts) error {
	logger := opts.Logger
	logger.Debug("apply nginx", zap.Any("opts", opts))

	// start nginx container
	if err := ensureNginxContainer(ctx, dockerClient, opts); err != nil {
		return errcode.TODO.Wrap(err)
	}

	// generate nginx config
	containersInfo, err := pwcompose.GetContainersInfo(ctx, dockerClient)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}
	if opts.DomainSuffix == "local" {
		proxyNetworkIP := containersInfo.NginxContainer.NetworkSettings.Networks[pwcompose.ProxyNetworkName].IPAddress
		opts.DomainSuffix = proxyNetworkIP + ".xip.io"
	}
	config, err := genNginxConfig(apiInstances, containersInfo, opts)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	/*if logger.Check(zap.DebugLevel, "") != nil {
		fmt.Fprintln(os.Stderr, "config", godev.PrettyJSON(config))
	}*/

	nginxContainer := containersInfo.NginxContainer
	// configure custom 503 page
	custom, err := buildCustom503PageTar(config, logger)
	if err != nil {
		return errcode.ErrBuildCustom503Page.Wrap(err)
	}
	logger.Debug("copy 503.html into the container", zap.String("container-id", nginxContainer.ID))
	err = dockerClient.CopyToContainer(ctx, nginxContainer.ID, "/usr/share/nginx/html/", custom, types.CopyToContainerOptions{})
	if err != nil {
		return errcode.ErrCopyNginxConfigToContainer.Wrap(err)
	}
	// configure nginx binary
	buf, err := buildNginxConfigTar(config, logger)
	if err != nil {
		return errcode.ErrBuildNginxConfig.Wrap(err)
	}
	logger.Debug("copy nginx config into the container", zap.String("container-id", nginxContainer.ID))
	err = dockerClient.CopyToContainer(ctx, nginxContainer.ID, "/etc/nginx/", buf, types.CopyToContainerOptions{})
	if err != nil {
		return errcode.ErrCopyNginxConfigToContainer.Wrap(err)
	}
	args := []string{"nginx", "-t", "-c", "/etc/nginx/nginx.conf"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	err = nginxSendCommand(ctx, dockerClient, nginxContainer.ID, logger, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandNewConfigCheck.Wrap(err)
	}
	// new config hot reload
	args = []string{"nginx", "-s", "reload"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	err = nginxSendCommand(ctx, dockerClient, nginxContainer.ID, logger, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandReloadConfig.Wrap(err)
	}
	/*if logger.Check(zap.DebugLevel, "") != nil {
		for _, upstream := range config.Upstreams {
			fmt.Fprintf(os.Stderr, "- %s\n", upstream.Name)
			for _, hash := range upstream.Hashes {
				fmt.Fprintf(os.Stderr, "  - %s.%s\n", hash, opts.DomainSuffix)
			}

		}
	}*/

	return nil
}

func ensureNginxContainer(ctx context.Context, dockerClient *client.Client, opts Opts) error {
	logger := opts.Logger

	// check if nginx server is started
	nginxContainer, err := checkNginxContainer(ctx, dockerClient)
	if err != nil {
		return errcode.ErrCheckNginxContainer.Wrap(err)
	}

	// remove nginx container if forced
	if opts.ForceRecreate && nginxContainer != nil {
		logger.Debug("remove old nginx", zap.String("id", nginxContainer.ID))
		err := dockerClient.ContainerRemove(ctx, nginxContainer.ID, types.ContainerRemoveOptions{
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
	var running bool
	if nginxContainer == nil {
		logger.Debug("build nginx", zap.Any("opts", opts))
		nginxContainerID, err = buildNginxContainer(ctx, dockerClient, opts)
		if err != nil {
			return errcode.ErrBuildNginxContainer.Wrap(err)
		}
		running = false
	} else {
		nginxContainerID = nginxContainer.ID
		running = nginxContainer.State == "running"
	}

	// start nginx container if needed
	if !running {
		err = dockerClient.ContainerStart(ctx, nginxContainerID, types.ContainerStartOptions{})
		if err != nil {
			return errcode.ErrStartNginxContainer.Wrap(err)
		}
		logger.Info("started nginx", zap.String("id", nginxContainerID))
	}

	// connect nginx container to proxy network
	nginxContainer, err = checkNginxContainer(ctx, dockerClient)
	if err != nil {
		return errcode.ErrCheckNginxContainer.Wrap(err)
	}
	if _, found := nginxContainer.NetworkSettings.Networks[pwcompose.ProxyNetworkName]; !found {
		var proxyNetworkID string
		networkResources, err := dockerClient.NetworkList(ctx, types.NetworkListOptions{})
		if err != nil {
			return errcode.ErrDockerAPINetworkList.Wrap(err)
		}
		for _, networkResource := range networkResources {
			if networkResource.Name == pwcompose.ProxyNetworkName {
				proxyNetworkID = networkResource.ID
				break
			}
		}
		logger.Debug("connect nginx network", zap.String("nginx-id", nginxContainer.ID), zap.String("network-id", proxyNetworkID))
		err = dockerClient.NetworkConnect(ctx, proxyNetworkID, nginxContainer.ID, nil)
		if err != nil {
			return errcode.ErrNginxConnectNetwork.Wrap(err)
		}
	}
	return nil
}

func genNginxConfig(apiInstances *pwapi.AgentListInstances_Output, containersInfo *pwcompose.ContainersInfo, opts Opts) (*nginxConfig, error) {
	config := nginxConfig{
		Opts:      opts,
		Upstreams: map[string]nginxUpstream{},
	}

	// compute allowed users by instance
	allowedUsers := map[string][]int64{}
	for _, apiInstance := range apiInstances.GetInstances() {
		if apiInstance.Status == pwdb.ChallengeInstance_Disabled {
			continue
		}
		uniqueUsers := map[int64]bool{}
		for _, seasonChallenge := range apiInstance.GetFlavor().GetSeasonChallenges() {
			for _, subscription := range seasonChallenge.GetActiveSubscriptions() {
				for _, member := range subscription.GetTeam().GetMembers() {
					uniqueUsers[member.UserID] = true
				}
			}
		}
		instanceID := fmt.Sprintf("%d", apiInstance.ID)
		allowedUsers[instanceID] = make([]int64, len(uniqueUsers))
		i := 0
		for user := range uniqueUsers {
			allowedUsers[instanceID][i] = user
			i++
		}
	}

	// compute upstreams
	for _, flavor := range containersInfo.RunningFlavors {
		for _, container := range flavor.Containers {
			for idx, port := range container.Ports {
				if port.PublicPort != 0 {
					upstream := nginxUpstream{
						Name:         fmt.Sprintf("%s.%d", container.Names[0][1:], idx),
						InstanceID:   flavor.InstanceKey,
						AllowedUsers: allowedUsers[flavor.InstanceKey],
						Host:         container.NetworkSettings.Networks[pwcompose.ProxyNetworkName].IPAddress,
						Port:         strconv.Itoa(int(port.PrivatePort)),
					}
					config.Upstreams[upstream.Name] = upstream
				}
			}
		}
	}

	for idx, upstream := range config.Upstreams {
		upstream.Hashes = make([]string, len(upstream.AllowedUsers))
		for j, userID := range upstream.AllowedUsers {
			hash, err := pwdb.ChallengeInstancePrefixHash(upstream.InstanceID, userID, opts.AuthSalt)
			if err != nil {
				return nil, errcode.ErrGeneratePrefixHash.Wrap(err)
			}
			upstream.Hashes[j] = hash
		}
		config.Upstreams[idx] = upstream
	}

	return &config, nil
}

func buildNginxConfigTar(config *nginxConfig, logger *zap.Logger) (*bytes.Buffer, error) {
	configTemplate, err := template.New("nginx-config").Parse(nginxConfigTemplate)
	if err != nil {
		return nil, errcode.ErrParsingTemplate.Wrap(err)
	}
	var configBuf bytes.Buffer
	err = configTemplate.Execute(&configBuf, config)
	if err != nil {
		return nil, errcode.ErrExecuteTemplate.Wrap(err)
	}
	configBytes := configBuf.Bytes()

	logger.Debug("nginx-config", zap.Int("config-length", len(configBytes)))
	/* if logger.Check(zap.DebugLevel, "") != nil {
		fmt.Fprintln(os.Stderr, string(configBytes))
	}*/

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: "nginx.conf",
		Mode: 0o755,
		Size: int64(len(configBytes)),
	})
	if err != nil {
		return nil, errcode.ErrWriteConfigFileHeader.Wrap(err)
	}

	if _, err := tw.Write(configBytes); err != nil {
		return nil, errcode.ErrWriteConfigFile.Wrap(err)
	}

	err = tw.Close()
	if err != nil {
		return nil, errcode.ErrCloseTarWriter.Wrap(err)
	}

	return &buf, nil
}

func buildCustom503PageTar(config *nginxConfig, logger *zap.Logger) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err := tw.WriteHeader(&tar.Header{
		Name: "503.html",
		Mode: 0o755,
		Size: int64(len(Template_503)),
	})
	if err != nil {
		return nil, errcode.ErrWriteCustom404FileHeader.Wrap(err)
	}

	if _, err := tw.Write([]byte(Template_503)); err != nil {
		return nil, errcode.ErrWriteCustom503File.Wrap(err)
	}

	err = tw.Close()
	if err != nil {
		return nil, errcode.ErrCloseTarWriter.Wrap(err)
	}

	return &buf, nil
}

func buildNginxContainer(ctx context.Context, cli *client.Client, opts Opts) (string, error) {
	logger := opts.Logger

	out, err := cli.ImagePull(ctx, opts.NginxDockerImage, types.ImagePullOptions{})
	if err != nil {
		return "", errcode.ErrDockerAPIImagePull.Wrap(err)
	}
	if logger.Check(zap.DebugLevel, "") != nil {
		_, err = io.Copy(os.Stderr, out)
		if err != nil {
			return "", errcode.TODO.Wrap(err)
		}
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
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
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

	if len(stderr) > 0 {
		logger.Debug("exec finished with stderr",
			zap.Strings("args", args),
			zap.String("stderr", string(stderr)),
		)
	}
	logger.Debug("exec finished",
		zap.Strings("args", args),
		zap.String("stdout", string(stdout)),
		zap.Int("exit-code", inspect.ExitCode),
		zap.Bool("running", inspect.Running),
	)

	if inspect.ExitCode != 0 {
		return errcode.ErrDockerAPIExitCode
	}

	return nil
}

type nginxConfig struct {
	Opts      Opts
	Upstreams map[string]nginxUpstream
}

type nginxUpstream struct {
	InstanceID   string
	Name         string
	Host         string
	Port         string
	Hashes       []string
	AllowedUsers []int64
}

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
	error_page 503 /503.html;
	location = /503.html {
			root /usr/share/nginx/html;
			internal;
	}
	location / {
			return 503;
	}
  }

  {{range .Upstreams -}}
  upstream upstream_{{.Name}} { server {{.Host}}:{{.Port}}; }
  server {
    listen      80;
    server_name moderator-{{.Name}}.{{$root.Opts.DomainSuffixWithoutPort}};
    access_log  /proc/self/fd/1;
    error_log   /proc/self/fd/2;
    # FIXME: add auth
    location = /robots.txt {
       add_header Content-Type text/plain;
       return 200 "User-agent: *\nDisallow: /\n";
    }
    location / {
      proxy_http_version 1.1;
      proxy_set_header Host                $http_host;
      proxy_set_header X-Real-IP           $remote_addr;
      proxy_set_header X-Forwarded-For     $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto   $scheme;
      proxy_set_header X-Frame-Options     SAMEORIGIN;
      proxy_set_header X-Pathwar-Mode      "moderator";
      proxy_set_header Upgrade             $http_upgrade;
      proxy_set_header Connection          "upgrade";
      proxy_pass http://upstream_{{.Name}};
    }
  }
  {{- if not (eq (len .Hashes) 0) }}
  server {
    listen      80;
    server_name{{range .Hashes}} {{.}}.{{$root.Opts.DomainSuffixWithoutPort}}{{end}};
    access_log  /proc/self/fd/1;
    error_log   /proc/self/fd/2;
    location = /robots.txt {
       add_header Content-Type text/plain;
       return 200 "User-agent: *\nDisallow: /\n";
    }
    location / {
      proxy_http_version  1.1;
      proxy_set_header Host                $http_host;
      proxy_set_header X-Real-IP           $remote_addr;
      proxy_set_header X-Forwarded-For     $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto   $scheme;
      proxy_set_header X-Frame-Options     SAMEORIGIN;
      proxy_set_header X-Pathwar-Mode      "authenticated";
      proxy_set_header Upgrade             $http_upgrade;
      proxy_set_header Connection          "upgrade";
      proxy_pass http://upstream_{{.Name}};
    }
  }
  {{end}}
  {{end -}}
}
`

const Template_503 = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Challenge is launching</title>
</head>
<body>
<h1 style="font-size: 1.5rem; font-weight: 400; line-height: 2.5rem; color: rgb(0, 129, 255); font-family: Bungee, cursive;">Challenge is lauching</h1>
<p style="font-size: 1.25rem; font-weight: bold; line-height: 1.1; color: #072A44; font-family: Barlow, sans-serif;">The challenge is launching, please try again in a few seconds.</p>
</body>
</html>
`
