package pwcompose

// https://github.com/digibib/docker-compose-dot/blob/master/docker-compose-dot.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

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
}

type dabfile struct {
	Services map[string]dabservice
}

type dabservice struct {
	Image string
}

func Prepare(path string, noPush bool, logger *zap.Logger) error {
	logger.Debug("prepare", zap.Bool("no-push", noPush), zap.String("path", path))

	path = filepath.Clean(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Fatal("wrong path: ", zap.Error(err))
	}

	challengeName := filepath.Base(path)

	// parse yaml
	composeData, err := ioutil.ReadFile(path + "/docker-compose.yml")
	if err != nil {
		logger.Fatal("docker-compose.yml file not found: ", zap.Error(err))
	}

	composeStruct := config{}
	err = yaml.Unmarshal(composeData, &composeStruct)
	if err != nil {
		logger.Fatal("error occured while parsing yaml: ", zap.Error(err))
	}

	// check yaml and add image name if not defined
	for name, service := range composeStruct.Services {
		if len(service.Image) == 0 {
			service.Image = "zarakii/" + challengeName + "." + name + ":latest"
			composeStruct.Services[name] = service
		}
	}

	// Rename original docker-compose file before creating the temporary one
	err = os.Rename(path+"/docker-compose.yml", path+"/docker-compose.yml.orig")
	if err != nil {
		logger.Fatal("file renaming error: ", zap.Error(err))
	}

	// create tmp docker-compose file
	tmpData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}

	tmpFile, err := os.Create(path + "/docker-compose.yml")
	if err != nil {
		logger.Fatal("temp file creation error: ", zap.Error(err))
	}
	tmpFile.Write(tmpData)
	tmpFile.Sync()

	// build and push images to dockerhub (don't forget to setup your credentials just type : "docker login" in bash)
	cmd := exec.Command("docker-compose", "build")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		logger.Fatal("docker-compose build error: ", zap.Error(err))
	} else {
		logger.Info("docker-compose build output: ", zap.String("bundle exec", string(out)))
	}

	cmd = exec.Command("docker-compose", "bundle", "--push-images")
	cmd.Dir = path
	out, err = cmd.Output()
	if err != nil {
		logger.Fatal("docker-compose bundle error: ", zap.Error(err))
	} else {
		logger.Info("bundle output: ", zap.String("bundle exec", string(out)))
	}

	// rename original docker-compose file and erase tmp one
	err = os.Rename(path+"/docker-compose.yml.orig", path+"/docker-compose.yml")
	if err != nil {
		logger.Fatal("file renaming error: ", zap.Error(err))
	}

	// parse json from .dab file
	composeDabfileJSON := dabfile{}
	composeDabfile, err := ioutil.ReadFile(path + "/" + challengeName + ".dab")
	if err != nil {
		logger.Fatal(challengeName+".dab file not found: ", zap.Error(err))
	}
	err = json.Unmarshal(composeDabfile, &composeDabfileJSON)

	// parse .dab file from docker-compose bundle
	composeData, err = yaml.Marshal(&composeStruct)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}

	// replace images from original docker-compose file with the one pushed to dockerhub
	for name, service := range composeStruct.Services {
		service.Image = composeDabfileJSON.Services[name].Image
		composeStruct.Services[name] = service
	}

	// cleanup
	err = os.Remove(path + "/" + challengeName + ".dab")
	if err != nil {
		logger.Fatal("couldn't remove "+challengeName+".dab file: ", zap.Error(err))
	}

	// print yaml
	// create tmp docker-compose file
	finalData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		logger.Fatal("error: ", zap.Error(err))
	}
	logger.Info("final docker-compose struct: ", zap.String("docker-compose.yml", string(finalData)))

	// return final docker-compose that can be used to deploy a challenge
	return fmt.Errorf("not implemented")
}

func Up(preparedCompose string, instanceKey string, logger *zap.Logger) error {
	logger.Debug(
		"up",
		zap.String("compose", preparedCompose),
		zap.String("instance-key", instanceKey),
	)
	// parse compose labels to get flavor info
	// generate instanceID with flavor+instanceKey
	// exec docker-compose up -d with params
	// print instanceID
	// later: print instance passphrases
	return fmt.Errorf("not implemented")
}

func Down(ids []string, logger *zap.Logger) error {
	logger.Debug("down", zap.Strings("ids", ids))
	// id can be an instance_id or a flavor_id
	return fmt.Errorf("not implemented")
}

func PS(depth int, logger *zap.Logger) error {
	logger.Debug("ps", zap.Int("depth", depth))
	return fmt.Errorf("not implemented")
}
