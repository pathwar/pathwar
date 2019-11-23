package pwhypervisor

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

func Nginx(ctx context.Context, config NginxConfig, cli *client.Client, logger *zap.Logger) error {
	logger.Debug("nginx", zap.Any("config", config))

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

	return fmt.Errorf("not implemented")
}

type NginxConfig struct {
	DomainSuffix      string             // .127.0.0.1.xip.io, .fr1.pathwar.pw, ...
	HTTPBind          string             // :8000, 0.0.0.0:80, ...
	ModeratorPassword string             // s3cur3
	Salt              string             // s3cur3-t0o
	AllowedUsers      map[string][]int64 // map[INSTANCE_ID][]USER_ID, map[42][]string{4242, 4343}
}
