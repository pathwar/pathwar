package pwcompose

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func Prepare(challengeDir string, prefix string, noPush bool, logger *zap.Logger) error {
	logger.Debug("prepare", zap.Bool("no-push", noPush), zap.String("challenge-dir", challengeDir), zap.String("prefix", prefix))

	cleanPath, err := filepath.Abs(filepath.Clean(challengeDir))
	if err != nil {
		return fmt.Errorf("get challenge dir: %w", err)
	}

	var (
		challengeName   = filepath.Base(cleanPath)
		origComposePath = path.Join(cleanPath, "docker-compose.yml")
		tmpComposePath  = path.Join(cleanPath, "docker-compose.tmp.yml")
		dabPath         = path.Join(cleanPath, challengeName+".dab")
	)

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return fmt.Errorf("challenge dir does not exist: %w", err)
	}

	// parse docker-compose.yml file
	composeData, err := ioutil.ReadFile(origComposePath)
	if err != nil {
		return fmt.Errorf("read docker-compose.yml: %w", err)
	}

	composeStruct := config{}
	err = yaml.Unmarshal(composeData, &composeStruct)
	if err != nil {
		return fmt.Errorf("parse docker-compose.yml: %w", err)
	}

	// check yaml and add image name if not defined
	for name, service := range composeStruct.Services {
		if len(service.Image) == 0 {
			service.Image = prefix + challengeName + "." + name
			composeStruct.Services[name] = service
		}
	}

	// create tmp docker-compose file
	tmpData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	tmpFile, err := os.Create(tmpComposePath)
	if err != nil {
		return fmt.Errorf("create tmp compose file: %w", err)
	}
	defer func() {
		if err = os.Remove(tmpComposePath); err != nil {
			logger.Warn("rm tmp compose file", zap.Error(err))
		}
	}()
	_, err = tmpFile.Write(tmpData)
	if err != nil {
		return fmt.Errorf("write tmp compose file: %w", err)
	}
	err = tmpFile.Sync()
	if err != nil {
		return fmt.Errorf("sync tmp compose file: %w", err)
	}
	tmpFile.Close()

	// build and push images to dockerhub (don't forget to setup your credentials just type : "docker login" in bash)
	logger.Debug("docker-compose", zap.String("-f", tmpComposePath), zap.String("action", "build"))
	cmd := exec.Command("docker-compose", "-f", tmpComposePath, "build")
	cmd.Dir = cleanPath
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker-compose build: %w", err)
	}

	logger.Debug("docker-compose", zap.String("-f", tmpComposePath), zap.String("action", "bundle"))
	pushImages := ""
	if !noPush {
		pushImages = "--push-images"
	}
	cmd = exec.Command("docker-compose", "-f", tmpComposePath, "bundle", pushImages)
	cmd.Dir = cleanPath
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker-compose bundle: %w", err)
	}
	defer func() {
		if err = os.Remove(dabPath); err != nil {
			logger.Warn("rm dab file", zap.Error(err))
		}
	}()

	// parse json from .dab file
	composeDabfileJSON := dabfile{}
	composeDabfile, err := ioutil.ReadFile(dabPath)
	if err != nil {
		return fmt.Errorf("read dab file: %w", err)
	}
	if err = json.Unmarshal(composeDabfile, &composeDabfileJSON); err != nil {
		return fmt.Errorf("parse dab: %w", err)
	}

	// replace images from original docker-compose file with the one pushed to dockerhub
	for name, service := range composeStruct.Services {
		service.Image = composeDabfileJSON.Services[name].Image
		composeStruct.Services[name] = service
	}

	// print yaml
	finalData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		return fmt.Errorf("marshal compose file: %w", err)
	}
	fmt.Println(string(finalData))

	return nil
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
