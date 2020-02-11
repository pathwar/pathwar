package pwcompose

import (
	"archive/tar"
	"bytes"
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
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"pathwar.land/go/internal/randstring"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwinit"
)

const (
	labelPrefix           = "land.pathwar.compose."
	serviceNameLabel      = labelPrefix + "service-name"
	serviceOrigin         = labelPrefix + "origin"
	challengeNameLabel    = labelPrefix + "challenge-name"
	challengeVersionLabel = labelPrefix + "challenge-version"
	instanceKeyLabel      = labelPrefix + "instance-key"
)

const (
	NginxContainerName = "pathwar-agent-nginx"
	ProxyNetworkName   = "pathwar-proxy-network"
)

func Prepare(challengeDir string, prefix string, noPush bool, version string, logger *zap.Logger) (string, error) {
	logger.Debug("prepare", zap.Bool("no-push", noPush), zap.String("challenge-dir", challengeDir), zap.String("prefix", prefix), zap.String("version", version))

	cleanPath, err := filepath.Abs(filepath.Clean(challengeDir))
	if err != nil {
		return "", errcode.ErrComposeInvalidPath.Wrap(err)
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
		return "", errcode.ErrComposeDirectoryNotFound.Wrap(err)
	}

	// parse docker-compose.yml file
	composeData, err := ioutil.ReadFile(origComposePath)
	if err != nil {
		return "", errcode.ErrComposeReadConfig.Wrap(err)
	}

	// check for error in docker-compose file
	args := []string{"-f", origComposePath, "config", "-q"}
	logger.Debug("docker-compose", zap.Strings("args", args))
	cmd := exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", errcode.ErrComposeInvalidConfig.Wrap(err)
	}

	composeStruct := config{}
	err = yaml.Unmarshal(composeData, &composeStruct)
	if err != nil {
		return "", errcode.ErrComposeInvalidConfig.Wrap(err)
	}

	// check yaml and add image name if not defined
	for name, service := range composeStruct.Services {
		if service.Labels == nil {
			service.Labels = map[string]string{}
		}
		if service.Image == "" {
			if !noPush {
				service.Image = prefix + challengeName + ":" + name
				service.Labels[serviceOrigin] = "was-built"
			} else {
				service.Build = path.Join(cleanPath, service.Build)
				service.Labels[serviceOrigin] = "was-built-dev"
			}
		} else {
			service.Labels[serviceOrigin] = "was-pulled"
		}
		service.Labels[challengeNameLabel] = challengeName
		service.Labels[serviceNameLabel] = name
		service.Labels[challengeVersionLabel] = version
		composeStruct.Services[name] = service
	}

	// create tmp docker-compose file
	tmpData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		return "", errcode.ErrComposeMarshalConfig.Wrap(err)
	}
	tmpFile, err := os.Create(tmpComposePath)
	if err != nil {
		return "", errcode.ErrComposeCreateTempFile.Wrap(err)
	}
	defer func() {
		if err = os.Remove(tmpComposePath); err != nil {
			logger.Warn("rm tmp compose file", zap.Error(err), zap.String("path", tmpComposePath))
		}
	}()
	_, err = tmpFile.Write(tmpData)
	if err != nil {
		return "", errcode.ErrComposeWriteTempFile.Wrap(err)
	}
	tmpFile.Close()

	if !noPush {
		// build and push images to dockerhub (don't forget to setup your credentials just type : "docker login" in bash)
		args = []string{"-f", tmpComposePath, "build"}
		logger.Debug("docker-compose", zap.Strings("args", args))
		cmd = exec.Command("docker-compose", args...)
		cmd.Dir = cleanPath
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return "", errcode.ErrComposeBuild.Wrap(err)
		}

		args := []string{"-f", tmpComposePath, "bundle", "--push-images"}
		logger.Debug("docker-compose", zap.Strings("args", args))
		cmd = exec.Command("docker-compose", args...)
		cmd.Dir = cleanPath
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return "", errcode.ErrComposeBundle.Wrap(err)
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
			return "", errcode.ErrComposeReadDab.Wrap(err)
		}
		if err = json.Unmarshal(composeDabfile, &composeDabfileJSON); err != nil {
			return "", errcode.ErrComposeParseDab.Wrap(err)
		}

		// replace images from original docker-compose file with the one pushed to dockerhub
		for name, service := range composeStruct.Services {
			service.Image = composeDabfileJSON.Services[name].Image
			service.Build = "" // ensure service only has an `image:` without a `build:`
			composeStruct.Services[name] = service
		}
	}

	// print yaml
	finalData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		return "", errcode.ErrComposeMarshalConfig.Wrap(err)
	}

	return string(finalData), nil
}

func Up(
	ctx context.Context,
	preparedCompose string,
	instanceKey string,
	forceRecreate bool,
	pwinitConfig *pwinit.InitConfig,
	cli *client.Client,
	logger *zap.Logger,
) error {
	logger.Debug("up", zap.String("compose", preparedCompose), zap.String("instance-key", instanceKey))

	// parse prepared compose yaml
	preparedComposeStruct := config{}
	err := yaml.Unmarshal([]byte(preparedCompose), &preparedComposeStruct)
	if err != nil {
		return errcode.ErrComposeParseConfig.Wrap(err)
	}

	// generate instanceIDs and set them as container_name
	var challengeID string
	for name, service := range preparedComposeStruct.Services {
		challengeName := service.Labels[challengeNameLabel]
		serviceName := service.Labels[serviceNameLabel]
		imageHash := "local"
		if strings.Contains(service.Image, "@sha256:") {
			imageHash = strings.Split(service.Image, "@sha256:")[1][:6]
		}
		service.ContainerName = fmt.Sprintf("%s.%s.%s.%s", challengeName, serviceName, imageHash, instanceKey)
		service.Restart = "unless-stopped"
		service.Labels[instanceKeyLabel] = instanceKey
		challengeID = challengeIDFormatted(service.Labels[challengeNameLabel], service.Labels[challengeVersionLabel])
		preparedComposeStruct.Services[name] = service
	}

	// down instances if force recreate
	if forceRecreate {
		err = Down(ctx, []string{challengeID}, false, false, false, cli, logger)
		if err != nil {
			return errcode.ErrComposeForceRecreateDown.Wrap(err)
		}
	}

	// create temp dir
	tmpDir, err := ioutil.TempDir("", "pwcompose")
	if err != nil {
		return errcode.ErrComposeCreateTempDir.Wrap(err)
	}
	defer func() {
		if err = os.RemoveAll(tmpDir); err != nil {
			logger.Warn("rm tmp dir", zap.Error(err))
		}
	}()
	tmpDirCompose := path.Join(tmpDir, challengeID)
	err = os.MkdirAll(tmpDirCompose, os.ModePerm)
	if err != nil {
		return errcode.ErrComposeCreateTempDir.Wrap(err)
	}

	// generate tmp path
	tmpPreparedComposePath := filepath.Join(tmpDirCompose, "docker-compose.yml")

	// create tmp docker-compose file
	err = updateDockerComposeTempFile(preparedComposeStruct, tmpPreparedComposePath)
	if err != nil {
		return errcode.ErrComposeUpdateTempFile.Wrap(err)
	}

	// create instances
	args := []string{"-f", tmpPreparedComposePath, "up", "--no-start"}
	logger.Debug("docker-compose", zap.Strings("args", args))
	cmd := exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Error("Error detected while creating containers, it's probably due to a conflict with previously created containers that share the same name. You should retry with --force-recreate flag")
		return errcode.ErrComposeRunCreate.Wrap(err)
	}

	// update entrypoints to run pwinit
	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}
	for _, instance := range containersInfo.RunningInstances {
		instanceChallengeID := challengeIDFormatted(instance.Labels[challengeNameLabel], instance.Labels[challengeVersionLabel])
		if challengeID == instanceChallengeID {
			// update entrypoints to run pwinit first
			imageInspect, _, err := cli.ImageInspectWithRaw(ctx, instance.ImageID)
			if err != nil {
				return errcode.ErrDockerAPIImageInspect.Wrap(err)
			}
			for name, service := range preparedComposeStruct.Services {
				if name != instance.Labels[serviceNameLabel] {
					continue
				}
				// find service from compose file of current instance
				entrypoint := []string{}
				if len(imageInspect.Config.Entrypoint) > 0 {
					entrypoint = imageInspect.Config.Entrypoint
				}
				if len(service.Entrypoint) > 0 {
					entrypoint = service.Entrypoint
				}
				command := []string{}
				if len(imageInspect.Config.Cmd) > 0 {
					command = imageInspect.Config.Cmd
				}
				if len(service.Command) > 0 {
					command = service.Command
				}
				service.Entrypoint = strslice.StrSlice{"/bin/pwinit", "entrypoint"}
				service.Command = append(entrypoint, command...)
				preparedComposeStruct.Services[name] = service
			}
		}
	}

	// update tmp docker-compose file with new entrypoints
	err = updateDockerComposeTempFile(preparedComposeStruct, tmpPreparedComposePath)
	if err != nil {
		return errcode.ErrComposeUpdateTempFile.Wrap(err)
	}

	// build definitive instances
	args = []string{"-f", tmpPreparedComposePath, "up", "--no-start"}
	logger.Debug("docker-compose", zap.Strings("args", args))
	cmd = exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logger.Error("Error detected while creating containers, it's probably due to a conflict with previously created containers that share the same name. You should retry with --force-recreate flag")
		return errcode.ErrComposeRunCreate.Wrap(err)
	}

	// copy pathwar binary inside all containers
	containersInfo, err = GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	for _, instance := range containersInfo.RunningInstances {
		if challengeID == challengeIDFormatted(instance.Labels[challengeNameLabel], instance.Labels[challengeVersionLabel]) {
			if pwinitConfig == nil {
				pwinitConfig = &pwinit.InitConfig{
					Passphrases: []string{
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
						fmt.Sprintf("dev-%s", randstring.RandString(10)),
					},
				}
			}
			buf, err := buildPWInitTar(*pwinitConfig)
			if err != nil {
				return errcode.ErrCopyPWInitToContainer.Wrap(err)
			}
			logger.Debug("copy pwinit into the container", zap.String("container-id", instance.ID))
			err = cli.CopyToContainer(ctx, instance.ID, "/", buf, types.CopyToContainerOptions{})
			if err != nil {
				return errcode.ErrCopyPWInitToContainer.Wrap(err)
			}
		}
	}

	// start instances
	args = []string{"-f", tmpPreparedComposePath, "up", "-d"}
	logger.Debug("docker-compose", zap.Strings("args", args))
	cmd = exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errcode.ErrComposeRunUp.Wrap(err)
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
	withNginx bool,
	cli *client.Client,
	logger *zap.Logger,
) error {
	logger.Debug("down", zap.Strings("ids", ids), zap.Bool("rmi", removeImages), zap.Bool("rm -v", removeVolumes), zap.Bool("with-nginx", withNginx))

	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	removalLists := dockerRemovalLists{
		containersToRemove: []string{},
		imagesToRemove:     []string{},
	}

	if withNginx && containersInfo.NginxProxyInstance.ID != "" {
		removalLists = updateDockerRemovalLists(removalLists, containersInfo.NginxProxyInstance, removeImages)
	}

	if len(ids) == 0 {
		for _, container := range containersInfo.RunningInstances {
			removalLists = updateDockerRemovalLists(removalLists, container, removeImages)
		}
	}

	for _, id := range ids {
		for _, flavor := range containersInfo.RunningFlavors {
			if id == flavor.Name || id == challengeIDFormatted(flavor.Name, flavor.Version) {
				for _, instance := range flavor.Instances {
					removalLists = updateDockerRemovalLists(removalLists, instance, removeImages)
				}
			}
		}
		for _, container := range containersInfo.RunningInstances {
			if id == container.ID || id == container.ID[0:7] {
				removalLists = updateDockerRemovalLists(removalLists, container, removeImages)
			}
		}
	}

	for _, instanceID := range removalLists.containersToRemove {
		err := cli.ContainerRemove(ctx, instanceID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: removeVolumes,
		})
		if err != nil {
			return errcode.ErrDockerAPIContainerRemove.Wrap(err)
		}
		fmt.Println("removed container " + instanceID)
	}

	for _, imageID := range removalLists.imagesToRemove {
		_, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		})
		if err != nil {
			return errcode.ErrDockerAPIImageRemove.Wrap(err)
		}
		fmt.Println("removed image " + imageID)
	}

	if withNginx && containersInfo.NginxProxyNetwork.ID != "" {
		err = cli.NetworkRemove(ctx, containersInfo.NginxProxyNetwork.ID)
		if err != nil {
			return errcode.ErrDockerAPINetworkRemove.Wrap(err)
		}
		fmt.Println("removed proxy network " + containersInfo.NginxProxyNetwork.ID)
	}

	return nil
}

func updateDockerRemovalLists(removalLists dockerRemovalLists, container types.Container, removeImages bool) dockerRemovalLists {
	removalLists.containersToRemove = append(removalLists.containersToRemove, container.ID)
	if removeImages {
		removalLists.imagesToRemove = append(removalLists.imagesToRemove, container.ImageID)
	}
	return removalLists
}

func PS(ctx context.Context, depth int, cli *client.Client, logger *zap.Logger) error {
	logger.Debug("ps", zap.Int("depth", depth))

	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "CHALLENGE", "SVC", "PORTS", "STATUS", "CREATED"})

	for _, flavor := range containersInfo.RunningFlavors {
		for uid, container := range flavor.Instances {

			ports := []string{}
			for _, port := range container.Ports {
				if port.PublicPort != 0 {
					ports = append(ports, strconv.Itoa(int(port.PublicPort)))
				}
			}

			table.Append([]string{
				uid[:7],
				challengeIDFormatted(flavor.Name, flavor.Version),
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

func buildPWInitTar(config pwinit.InitConfig) (*bytes.Buffer, error) {
	var pwInitBuf []byte
	pwInitBuf, err := pwinit.Binary()
	if err != nil {
		return nil, errcode.ErrGetPWInitBinary.Wrap(err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// write pwinit binary into tar file
	err = tw.WriteHeader(&tar.Header{
		Name: "/bin/pwinit",
		Mode: 0755,
		Size: int64(len(pwInitBuf)),
	})
	if err != nil {
		return nil, errcode.ErrWritePWInitFileHeader.Wrap(err)
	}
	_, err = tw.Write(pwInitBuf)
	if err != nil {
		return nil, errcode.ErrWritePWInitFile.Wrap(err)
	}

	// write pwinit json config into tar file
	pwInitConfigJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errcode.ErrMarshalPWInitConfigFile.Wrap(err)
	}
	err = tw.WriteHeader(&tar.Header{
		Name: "/pwinit/config.json",
		Mode: 0755,
		Size: int64(len(pwInitConfigJSON)),
		// FIXME: chown it to container's default user
	})
	if err != nil {
		return nil, errcode.ErrWritePWInitConfigFileHeader.Wrap(err)
	}
	_, err = tw.Write(pwInitConfigJSON)
	if err != nil {
		return nil, errcode.ErrWritePWInitConfigFile.Wrap(err)
	}

	if err = tw.Close(); err != nil {
		return nil, errcode.ErrWritePWInitCloseTarWriter.Wrap(err)
	}

	return &buf, nil
}

func GetContainersInfo(ctx context.Context, cli *client.Client) (*ContainersInfo, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, errcode.ErrDockerAPIContainerList.Wrap(err)
	}

	containersInfo := ContainersInfo{
		RunningFlavors:   map[string]challengeFlavors{},
		RunningInstances: map[string]types.Container{},
	}

	for _, container := range containers {
		// find nginx proxy
		for _, name := range container.Names {
			if name[1:] == NginxContainerName {
				containersInfo.NginxProxyInstance = container
			}
		}
		// continue if container is not a challenge
		if _, pwcontainer := container.Labels[challengeNameLabel]; !pwcontainer {
			continue
		}
		// handle and sort challenge
		flavor := fmt.Sprintf(
			"%s:%s",
			container.Labels[challengeNameLabel],
			container.Labels[challengeVersionLabel],
		)
		if _, found := containersInfo.RunningFlavors[flavor]; !found {
			challengeFlavor := challengeFlavors{
				Instances: map[string]types.Container{},
			}
			challengeFlavor.Name = container.Labels[challengeNameLabel]
			challengeFlavor.Version = container.Labels[challengeVersionLabel]
			challengeFlavor.InstanceKey = container.Labels[instanceKeyLabel]
			containersInfo.RunningFlavors[flavor] = challengeFlavor
		}
		containersInfo.RunningFlavors[flavor].Instances[container.ID] = container
		containersInfo.RunningInstances[container.ID] = container
	}

	// find proxy network
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, errcode.ErrDockerAPINetworkList.Wrap(err)
	}
	for _, networkResource := range networks {
		if networkResource.Name == ProxyNetworkName {
			containersInfo.NginxProxyNetwork = networkResource
			break
		}
	}

	return &containersInfo, nil
}

func updateDockerComposeTempFile(preparedComposeStruct config, tmpPreparedComposePath string) error {
	// create tmp docker-compose file
	tmpData, err := yaml.Marshal(&preparedComposeStruct)
	if err != nil {
		return errcode.ErrComposeMarshalConfig.Wrap(err)
	}
	tmpFile, err := os.Create(tmpPreparedComposePath)
	if err != nil {
		return errcode.ErrComposeCreateTempFile.Wrap(err)
	}

	_, err = tmpFile.Write(tmpData)
	if err != nil {
		return errcode.ErrComposeWriteTempFile.Wrap(err)
	}
	err = tmpFile.Close()
	if err != nil {
		return errcode.ErrComposeCloseTempFile.Wrap(err)
	}
	return nil
}

func challengeIDFormatted(challengeNameLabel string, challengeVersionLabel string) string {
	return fmt.Sprintf("%s@%s", challengeNameLabel, challengeVersionLabel)
}
