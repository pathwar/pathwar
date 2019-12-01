package pwhypervisor

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/pkg/stdcopy"
	"go.uber.org/zap"
	"pathwar.land/go/internal/randstring"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwcompose"
)

const nginxContainerName = "pathwar-hypervisor-nginx"

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
					_ = port.PublicPort
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
				Name: "instance1",
				Host: "1.2.3.4",
				Port: "1337",
			},
			{
				Name: "instance2",
				Host: "2.3.4.5",
				Port: "1338",
				Hashes: []string{
					"123",
					"456",
				},
			},
		},
		Opts: opts,
	}
	buf, err := buildNginxConfigTar(configData)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	logger.Debug("copy nginx config into the container", zap.String("container-id", nginxContainerID))
	err = cli.CopyToContainer(ctx, nginxContainerID, "/etc/nginx/", buf, types.CopyToContainerOptions{})
	if err != nil {
		return errcode.ErrCopyNginxConfigToContainer.Wrap(err)
	}

	// check new nginx config
	//args := []string{"nginx", "-t"}
	args := []string{"head", "/etc/nginx/nginx.conf"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	err = nginxSendCommand(ctx, cli, nginxContainerID, logger, args...)
	if err != nil {
		return errcode.ErrNginxSendCommandNewConfigCheck.Wrap(err)
	}

	// new config hot reload
	args = []string{"nginx", "-s", "reload"}
	logger.Debug("send nginx command", zap.Strings("args", args))
	err = nginxSendCommand(ctx, cli, nginxContainerID, logger, args...)
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

func buildNginxConfigTar(data interface{}) (*bytes.Buffer, error) {
	configTemplate, err := template.New("nginx-config").Parse(nginxConfigTemplate)
	if err != nil {
		return nil, errcode.ErrParsingTemplate.Wrap(err)
	}
	var configBuf bytes.Buffer
	err = configTemplate.Execute(&configBuf, data)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	config := configBuf.Bytes()
	// fmt.Println(string(config))

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

// https://github.com/pathwar/pathwar/blob/v1.0.0/hypervisor/cmd_hypervisor_run.go
func buildNginxContainer(ctx context.Context, cli *client.Client, opts HypervisorOpts) (string, error) {
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
		}, nil, nginxContainerName)
	if err != nil {
		return "", errcode.ErrDockerAPIContainerCreate.Wrap(err)
	}

	return cont.ID, nil
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

	/*err = cli.ContainerExecStart(ctx, execRes.ID, types.ExecStartCheck{})
	if err != nil {
		return errcode.ErrDockerAPIContainerExecStart.Wrap(err)
	}*/

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
