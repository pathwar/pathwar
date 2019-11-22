package pwcompose

import (
	"github.com/docker/docker/api/types"
)

// https://github.com/digibib/docker-compose-dot/blob/master/docker-compose-dot.go
type config struct {
	Version  string
	Networks map[string]network
	Volumes  map[string]volume
	Services map[string]service
}

type network struct {
	Driver, External string            `yaml:",omitempty"`
	DriverOpts       map[string]string `yaml:"driver_opts,omitempty"`
}

type volume struct {
	Driver, External string            `yaml:",omitempty"`
	DriverOpts       map[string]string `yaml:"driver_opts,omitempty"`
}

type service struct {
	ContainerName                             string            `yaml:"container_name,omitempty"`
	Image                                     string            `yaml:",omitempty"`
	Networks, Ports, Expose, Volumes, Command []string          `yaml:",omitempty"`
	VolumesFrom                               []string          `yaml:"volumes_from,omitempty"`
	DependsOn                                 []string          `yaml:"depends_on,omitempty"`
	CapAdd                                    []string          `yaml:"cap_add,omitempty"`
	Build                                     string            `yaml:",omitempty"`
	Environment                               map[string]string `yaml:",omitempty"`
	Labels                                    map[string]string `yaml:"labels,omitempty"`
}

type dabfile struct {
	Services map[string]dabservice
}

type dabservice struct {
	Image string
}

type pathwarInfo struct {
	RunningFlavors   map[string]challengeFlavors
	RunningInstances map[string]types.Container
}

type challengeFlavors struct {
	Name      string
	Version   string
	Instances map[string]types.Container
}
