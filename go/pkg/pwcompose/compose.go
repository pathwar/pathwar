package pwcompose

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	labelPrefix           = "land.pathwar.compose."
	serviceNameLabel      = labelPrefix + "service-name"
	challengeNameLabel    = labelPrefix + "challenge-name"
	challengeVersionLabel = labelPrefix + "challenge-version"
	instanceKeyLabel      = labelPrefix + "instance-key"
)

func Prepare(challengeDir string, prefix string, noPush bool, version string, logger *zap.Logger) error {
	logger.Debug("prepare", zap.Bool("no-push", noPush), zap.String("challenge-dir", challengeDir), zap.String("prefix", prefix), zap.String("version", version))

	cleanPath, err := filepath.Abs(filepath.Clean(challengeDir))
	if err != nil {
		return fmt.Errorf("get challenge dir: %w", err)
	}

	if prefix[len(prefix)-1:] != "/" {
		prefix += "/"
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
		if service.Image == "" {
			service.Image = prefix + challengeName + ":" + name
		}
		if service.Labels == nil {
			service.Labels = map[string]string{}
		}
		service.Labels[challengeNameLabel] = challengeName
		service.Labels[serviceNameLabel] = name
		service.Labels[challengeVersionLabel] = version
		composeStruct.Services[name] = service
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

	cmdArgs := []string{"docker-compose", "-f", tmpComposePath, "bundle"}
	if !noPush {
		cmdArgs = append(cmdArgs, "--push-images")
	}
	logger.Debug("docker-compose", zap.Strings("args", cmdArgs[1:]))
	cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
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
		service.Build = "" // ensure service only has an `image:` without a `build:`
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

	// parse prepared compose yaml
	preparedComposeStruct := config{}
	err := yaml.Unmarshal([]byte(preparedCompose), &preparedComposeStruct)
	if err != nil {
		return fmt.Errorf("parse prepared compose: %w", err)
	}

	// generate instanceIDs and set them as container_name
	for name, service := range preparedComposeStruct.Services {
		challengeName := service.Labels[challengeNameLabel]
		serviceName := service.Labels[serviceNameLabel]
		imageHash := strings.Split(service.Image, "@sha256:")[1]
		service.ContainerName = fmt.Sprintf("%s.%s.%s.%s", challengeName, serviceName, imageHash[:6], instanceKey)
		service.Labels[instanceKeyLabel] = instanceKey
		preparedComposeStruct.Services[name] = service
	}

	tmpDir, err := ioutil.TempDir("", "pwcompose")
	if err != nil {
		return fmt.Errorf("temp dir creation: %w", err)
	}
	defer func() {
		if err = os.RemoveAll(tmpDir); err != nil {
			logger.Warn("rm tmp dir", zap.Error(err))
		}
	}()

	tmpPreparedComposePath := filepath.Join(tmpDir, "docker-compose.yml")

	// create tmp docker-compose file
	tmpData, err := yaml.Marshal(&preparedComposeStruct)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	tmpFile, err := os.Create(tmpPreparedComposePath)
	if err != nil {
		return fmt.Errorf("create tmp compose file: %w", err)
	}

	_, err = tmpFile.Write(tmpData)
	if err != nil {
		return fmt.Errorf("write tmp compose file: %w", err)
	}
	tmpFile.Close()

	// start instances
	logger.Debug("docker-compose", zap.String("action", "up"))
	cmd := exec.Command("docker-compose", "-f", tmpPreparedComposePath, "up", "-d")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker-compose up -d: %w", err)
	}

	// print instanceIDs
	for _, service := range preparedComposeStruct.Services {
		fmt.Println(service.ContainerName)
	}

	return nil
}

func Down(
	ctx context.Context,
	ids []string,
	removeImages bool,
	removeVolumes bool,
	cli *client.Client,
	logger *zap.Logger,
) error {
	logger.Debug("down", zap.Strings("ids", ids), zap.Bool("rmi", removeImages), zap.Bool("rm -v", removeVolumes))

	pwInfo, err := getPathwarInfo(ctx, cli)
	if err != nil {
		return fmt.Errorf("error retrieving pathwar containers info: %w", err)
	}

	var (
		containersToRemove []string
		imagesToRemove     []string
	)

	for _, id := range ids {
		for _, flavor := range pwInfo.RunningFlavors {
			if id == flavor.Name || id == flavor.Name+":"+flavor.Version {
				for _, instance := range flavor.Instances {
					containersToRemove = append(containersToRemove, instance.ID)
					if removeImages == true {
						imagesToRemove = append(imagesToRemove, instance.ImageID)
					}
				}
			}
		}
		for _, container := range pwInfo.RunningInstances {
			if id == container.ID || id == container.ID[0:6] {
				containersToRemove = append(containersToRemove, container.ID)
				if removeImages == true {
					imagesToRemove = append(imagesToRemove, container.ImageID)
				}
			}
		}
	}

	for _, instanceID := range containersToRemove {
		err := cli.ContainerRemove(ctx, instanceID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: removeVolumes,
		})
		if err != nil {
			return fmt.Errorf("error while removing container: %w", err)
		}
		fmt.Println("removed container " + instanceID)
	}

	for _, imageID := range imagesToRemove {
		_, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			return fmt.Errorf("error while removing container: %w", err)
		}
		fmt.Println("removed image " + imageID)
	}

	return nil
}

func PS(ctx context.Context, depth int, cli *client.Client, logger *zap.Logger) error {
	logger.Debug("ps", zap.Int("depth", depth))

	pwInfo, err := getPathwarInfo(ctx, cli)
	if err != nil {
		return fmt.Errorf("retrieve containers info: %w", err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "CHALLENGE", "SVC", "PORTS", "STATUS", "CREATED"})

	for _, flavor := range pwInfo.RunningFlavors {
		for uid, container := range flavor.Instances {

			ports := []string{}
			for _, port := range container.Ports {
				if port.PublicPort != 0 {
					ports = append(ports, strconv.Itoa(int(port.PublicPort)))
				}
			}
			fmt.Println("AAA", ports, "BBB")

			table.Append([]string{
				uid[:7],
				fmt.Sprintf("%s@%s", flavor.Name, flavor.Version),
				container.Labels[serviceNameLabel],
				strings.Join(ports, ", "),
				strings.Replace(container.Status, "Up ", "", 1),
				strings.Replace(humanize.Time(time.Unix(container.Created, 0)), " ago", "", 1),
			})
		}
	}
	table.Render()
	return nil
}

func getPathwarInfo(ctx context.Context, cli *client.Client) (*pathwarInfo, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	pwInfo := pathwarInfo{
		RunningFlavors:   map[string]challengeFlavors{},
		RunningInstances: map[string]types.Container{},
	}

	for _, container := range containers {
		if _, pwcontainer := container.Labels[challengeNameLabel]; !pwcontainer {
			continue
		}
		flavor := fmt.Sprintf(
			"%s:%s",
			container.Labels[challengeNameLabel],
			container.Labels[challengeVersionLabel],
		)
		if _, found := pwInfo.RunningFlavors[flavor]; !found {
			challengeFlavor := challengeFlavors{
				Instances: map[string]types.Container{},
			}
			challengeFlavor.Name = container.Labels[challengeNameLabel]
			challengeFlavor.Version = container.Labels[challengeVersionLabel]
			pwInfo.RunningFlavors[flavor] = challengeFlavor
		}
		pwInfo.RunningFlavors[flavor].Instances[container.ID] = container
		pwInfo.RunningInstances[container.ID] = container
	}

	return &pwInfo, nil
}
