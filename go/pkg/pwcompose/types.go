package pwcompose

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

type config struct {
	Version  string
	Networks map[string]network
	Volumes  map[string]volume
	Services map[string]Service
	Pathwar  struct {
		Challenge pwdb.Challenge       `yaml:"challenge" json:"challenge"`
		Flavor    pwdb.ChallengeFlavor `yaml:"flavor" json:"flavor"`
	} `yaml:"x-pathwar" json:"pathwar"`
}

type network struct {
	Driver, External string            `yaml:",omitempty"`
	DriverOpts       map[string]string `yaml:"driver_opts,omitempty"`
}

type volume struct {
	Driver, External string            `yaml:",omitempty"`
	DriverOpts       map[string]string `yaml:"driver_opts,omitempty"`
}

type Service struct {
	ContainerName                             string            `yaml:"container_name,omitempty"`
	Image                                     string            `yaml:",omitempty"`
	Networks, Ports, Expose, Volumes, Command []string          `yaml:",omitempty"`
	VolumesFrom                               []string          `yaml:"volumes_from,omitempty"`
	DependsOn                                 []string          `yaml:"depends_on,omitempty"`
	CapAdd                                    []string          `yaml:"cap_add,omitempty"`
	Build                                     string            `yaml:",omitempty"`
	Entrypoint                                []string          `yaml:",omitempty"`
	Restart                                   string            `yaml:",omitempty"`
	Environment                               map[string]string `yaml:",omitempty"`
	Labels                                    map[string]string `yaml:"labels,omitempty"`
}

func (s Service) ChallengeID() string {
	return fmt.Sprintf("%s@%s", s.Labels[challengeNameLabel], s.Labels[challengeVersionLabel])
}

type dabfile struct {
	Services map[string]dabservice
}

type dabservice struct {
	Image string
}

type ContainersInfo struct {
	RunningFlavors    map[string]challengeFlavors
	RunningContainers map[string]container
	NginxContainer    container
	NginxNetwork      types.NetworkResource
}

type container types.Container

func (c container) ChallengeID() string {
	return fmt.Sprintf("%s@%s", c.Labels[challengeNameLabel], c.Labels[challengeVersionLabel])
}

func (c container) NeedsNginxProxy() bool {
	for _, port := range c.Ports {
		if port.PrivatePort != 0 {
			return true
		}
	}
	return false
}

type challengeFlavors struct {
	Name        string
	Version     string
	InstanceKey string
	Containers  map[string]container
}

func (cf challengeFlavors) ChallengeID() string {
	return fmt.Sprintf("%s@%s", cf.Name, cf.Version)
}
