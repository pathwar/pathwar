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
	"pathwar.land/pathwar/v2/go/internal/randstring"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

const (
	labelPrefix           = "land.pathwar.compose."
	serviceNameLabel      = labelPrefix + "service-name"
	serviceOrigin         = labelPrefix + "origin"
	challengeNameLabel    = labelPrefix + "challenge-name"
	challengeVersionLabel = labelPrefix + "challenge-version"
	InstanceKeyLabel      = labelPrefix + "instance-key"
)

const (
	NginxContainerName = "pathwar-agent-nginx"
	ProxyNetworkName   = "pathwar-proxy-network"
)

const (
	defaultDockerPrefix = "pathwar/"
)

type PrepareOpts struct {
	ChallengeDir string
	Prefix       string
	NoPush       bool
	Version      string
	Logger       *zap.Logger
}

func NewPrepareOpts() PrepareOpts {
	return PrepareOpts{
		ChallengeDir: ".",
		Prefix:       defaultDockerPrefix,
		Version:      "dev",
		NoPush:       false,
	}
}

func (opts *PrepareOpts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.ChallengeDir == "" {
		opts.ChallengeDir = "."
	}
}

func Prepare(opts PrepareOpts) (string, error) {
	opts.applyDefaults()
	opts.Logger.Debug("prepare", zap.Any("opts", opts))

	cleanPath, err := filepath.Abs(filepath.Clean(opts.ChallengeDir))
	if err != nil {
		return "", errcode.ErrComposeInvalidPath.Wrap(err)
	}

	if opts.Prefix[len(opts.Prefix)-1:] != "/" {
		opts.Prefix += "/"
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
	args := append(composeCliCommonArgs(origComposePath), "config", "-q")
	opts.Logger.Debug("docker-compose", zap.Strings("args", args))
	cmd := exec.Command("docker-compose", args...)
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
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
			if !opts.NoPush {
				service.Image = opts.Prefix + challengeName + ":" + name
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
		service.Labels[challengeVersionLabel] = opts.Version
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
			opts.Logger.Warn("rm tmp compose file", zap.Error(err), zap.String("path", tmpComposePath))
		}
	}()
	_, err = tmpFile.Write(tmpData)
	if err != nil {
		return "", errcode.ErrComposeWriteTempFile.Wrap(err)
	}
	tmpFile.Close()

	if !opts.NoPush {
		// build and push images to dockerhub (don't forget to setup your credentials just type : "docker login" in bash)
		args = append(composeCliCommonArgs(tmpComposePath), "build")
		opts.Logger.Debug("docker-compose", zap.Strings("args", args))
		cmd = exec.Command("docker-compose", args...)
		cmd.Dir = cleanPath
		if opts.Logger.Check(zap.DebugLevel, "") != nil {
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			return "", errcode.ErrComposeBuild.Wrap(err)
		}

		args = append(composeCliCommonArgs(tmpComposePath), "bundle", "--push-images")
		opts.Logger.Debug("docker-compose", zap.Strings("args", args))
		cmd = exec.Command("docker-compose", args...)
		cmd.Dir = cleanPath
		if opts.Logger.Check(zap.DebugLevel, "") != nil {
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			return "", errcode.ErrComposeBundle.Wrap(err)
		}
		defer func() {
			if err = os.Remove(dabPath); err != nil {
				opts.Logger.Warn("rm dab file", zap.Error(err))
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

type UpOpts struct {
	PreparedCompose string
	InstanceKey     string
	ForceRecreate   bool
	ProxyNetworkID  string
	PwinitConfig    *pwinit.InitConfig
	Logger          *zap.Logger
}

func NewUpOpts() UpOpts {
	return UpOpts{
		InstanceKey: "default",
	}
}

func (opts *UpOpts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.InstanceKey == "" {
		opts.InstanceKey = "default"
	}
	if opts.PwinitConfig == nil {
		opts.PwinitConfig = &pwinit.InitConfig{
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
}

// Up starts a prepared challenge
// nolint:gocyclo
func Up(ctx context.Context, cli *client.Client, opts UpOpts) (map[string]Service, error) {
	opts.applyDefaults()
	opts.Logger.Debug("up", zap.Any("opts", opts))

	// parse prepared compose yaml
	preparedComposeStruct := config{}
	err := yaml.Unmarshal([]byte(opts.PreparedCompose), &preparedComposeStruct)
	if err != nil {
		return nil, errcode.ErrComposeParseConfig.Wrap(err)
	}

	var challengeID string
	// generate instanceIDs and set them as container_name
	for name, service := range preparedComposeStruct.Services {
		challengeName := service.Labels[challengeNameLabel]
		serviceName := service.Labels[serviceNameLabel]
		imageHash := "local"
		if strings.Contains(service.Image, "@sha256:") {
			imageHash = strings.Split(service.Image, "@sha256:")[1][:6]
		}
		service.ContainerName = fmt.Sprintf("%s.%s.%s.%s", challengeName, serviceName, imageHash, opts.InstanceKey)
		service.Restart = "unless-stopped"
		service.Labels[InstanceKeyLabel] = opts.InstanceKey
		preparedComposeStruct.Services[name] = service
		if challengeID == "" {
			challengeID = service.ChallengeID()
		}
	}

	// down containers if force recreate
	if opts.ForceRecreate {
		err = Clean(ctx, cli, CleanOpts{Logger: opts.Logger})
		if err != nil {
			return nil, errcode.ErrComposeForceRecreateDown.Wrap(err)
		}
	}

	// create temp dir
	tmpDir, err := ioutil.TempDir("", "pwcompose")
	if err != nil {
		return nil, errcode.ErrComposeCreateTempDir.Wrap(err)
	}
	defer func() {
		if err = os.RemoveAll(tmpDir); err != nil {
			opts.Logger.Warn("rm tmp dir", zap.Error(err))
		}
	}()
	tmpDirCompose := path.Join(tmpDir, challengeID)
	err = os.MkdirAll(tmpDirCompose, os.ModePerm)
	if err != nil {
		return nil, errcode.ErrComposeCreateTempDir.Wrap(err)
	}

	// generate tmp path
	tmpPreparedComposePath := filepath.Join(tmpDirCompose, "docker-compose.yml")

	// create tmp docker-compose file
	err = updateDockerComposeTempFile(preparedComposeStruct, tmpPreparedComposePath)
	if err != nil {
		return nil, errcode.ErrComposeUpdateTempFile.Wrap(err)
	}

	// create containers
	args := append(composeCliCommonArgs(tmpPreparedComposePath), "up", "--no-start", "--quiet-pull")
	opts.Logger.Debug("docker-compose", zap.Strings("args", args))
	cmd := exec.Command("docker-compose", args...)
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		opts.Logger.Error("Error detected while creating containers, it's probably due to a conflict with previously created containers that share the same name. You should retry with --force-recreate flag")
		return nil, errcode.ErrComposeRunCreate.Wrap(err)
	}

	// update entrypoints to run pwinit
	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return nil, errcode.ErrComposeGetContainersInfo.Wrap(err)
	}
	for _, container := range containersInfo.RunningContainers {
		if challengeID == container.ChallengeID() {
			// update entrypoints to run pwinit first
			imageInspect, _, err := cli.ImageInspectWithRaw(ctx, container.ImageID)
			if err != nil {
				return nil, errcode.ErrDockerAPIImageInspect.Wrap(err)
			}
			for name, service := range preparedComposeStruct.Services {
				if name != container.Labels[serviceNameLabel] {
					continue
				}
				// find service from compose file of current container
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
				newCommand := entrypoint
				newCommand = append(newCommand, command...)
				service.Command = newCommand
				preparedComposeStruct.Services[name] = service
			}
		}
	}

	// update tmp docker-compose file with new entrypoints
	err = updateDockerComposeTempFile(preparedComposeStruct, tmpPreparedComposePath)
	if err != nil {
		return nil, errcode.ErrComposeUpdateTempFile.Wrap(err)
	}

	// if debug flag, display updated compose file
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		data, err := ioutil.ReadFile(tmpPreparedComposePath)
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}
		fmt.Println(string(data))
	}

	// build definitive containers
	args = append(composeCliCommonArgs(tmpPreparedComposePath), "up", "--no-start")
	opts.Logger.Debug("docker-compose", zap.Strings("args", args))
	cmd = exec.Command("docker-compose", args...)
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		opts.Logger.Error("Error detected while creating containers, it's probably due to a conflict with previously created containers that share the same name. You should retry with --force-recreate flag")
		return nil, errcode.ErrComposeRunCreate.Wrap(err)
	}

	// copy pathwar binary inside all containers
	containersInfo, err = GetContainersInfo(ctx, cli)
	if err != nil {
		return nil, errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	for _, container := range containersInfo.RunningContainers {
		if challengeID != container.ChallengeID() {
			continue
		}

		buf, err := buildPWInitTar(*opts.PwinitConfig)
		if err != nil {
			return nil, errcode.ErrCopyPWInitToContainer.Wrap(err)
		}
		opts.Logger.Debug("copy pwinit into the container", zap.String("container-id", container.ID))
		err = cli.CopyToContainer(ctx, container.ID, "/", buf, types.CopyToContainerOptions{})
		if err != nil {
			return nil, errcode.ErrCopyPWInitToContainer.Wrap(err)
		}
	}

	// start containers
	args = append(composeCliCommonArgs(tmpPreparedComposePath), "up", "-d")
	opts.Logger.Debug("docker-compose", zap.Strings("args", args))
	cmd = exec.Command("docker-compose", args...)
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		return nil, errcode.ErrComposeRunUp.Wrap(err)
	}

	// attach networks
	containersInfo, err = GetContainersInfo(ctx, cli) // this get containers info can be skipped if using the parsed compose file already in memory
	if err != nil {
		return nil, errcode.ErrComposeGetContainersInfo.Wrap(err)
	}
	for _, container := range containersInfo.RunningContainers {
		if challengeID != container.ChallengeID() {
			continue
		}
		if opts.ProxyNetworkID != "" && container.NeedsNginxProxy() {
			err = cli.NetworkConnect(ctx, opts.ProxyNetworkID, container.ID, nil)
			if err != nil {
				return nil, errcode.ErrContainerConnectNetwork.Wrap(err)
			}
		}
	}

	return preparedComposeStruct.Services, nil
}

func composeCliCommonArgs(path string) []string {
	return []string{"-f", path, "--no-ansi", "--log-level=ERROR"}
}

// Purge cleans up everything related to Pathwar (containers, volumes, images, networks)
func Purge(ctx context.Context, cli *client.Client, logger *zap.Logger) error {
	return Clean(ctx, cli, CleanOpts{
		RemoveImages:  true,
		RemoveVolumes: true,
		RemoveNginx:   true,
		Logger:        logger,
	})
}

// DownAll cleans up everything related to Pathwar except images (containers, volumes, networks)
func DownAll(ctx context.Context, cli *client.Client, logger *zap.Logger) error {
	return Clean(ctx, cli, CleanOpts{
		RemoveVolumes: true,
		RemoveNginx:   true,
		Logger:        logger,
	})
}

type CleanOpts struct {
	ContainerIDs  []string
	RemoveImages  bool
	RemoveVolumes bool
	RemoveNginx   bool
	Logger        *zap.Logger
}

func NewCleanOpts() CleanOpts {
	return CleanOpts{
		RemoveVolumes: true,
	}
}

func (opts *CleanOpts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
}

// Clean can cleanup specific containers, all the images, all the volumes, and the pathwar's nginx front-end
func Clean(ctx context.Context, cli *client.Client, opts CleanOpts) error {
	opts.applyDefaults()
	opts.Logger.Debug("down", zap.Any("opts", opts))

	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	toRemove := map[string]container{}

	if opts.RemoveNginx && containersInfo.NginxContainer.ID != "" {
		toRemove[containersInfo.NginxContainer.ID] = containersInfo.NginxContainer
	}

	if len(opts.ContainerIDs) == 0 { // all containers
		for _, container := range containersInfo.RunningContainers {
			toRemove[container.ID] = container
		}
	} else { // only specific ones
		for _, id := range opts.ContainerIDs {
			for _, flavor := range containersInfo.RunningFlavors {
				if id == flavor.Name || id == flavor.ChallengeID() {
					for _, container := range flavor.Containers {
						toRemove[container.ID] = container
					}
				}
			}
			for _, container := range containersInfo.RunningContainers {
				if id == container.ID || id == container.ID[0:7] {
					toRemove[container.ID] = container
				}
			}
		}
	}

	for _, container := range toRemove {
		err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: opts.RemoveVolumes,
		})
		if err != nil {
			return errcode.ErrDockerAPIContainerRemove.Wrap(err)
		}
		opts.Logger.Debug("container removed", zap.String("ID", container.ID))
		if opts.RemoveImages {
			_, err := cli.ImageRemove(ctx, container.ImageID, types.ImageRemoveOptions{
				Force:         false,
				PruneChildren: true,
			})
			if err != nil {
				return errcode.ErrDockerAPIImageRemove.Wrap(err)
			}
			opts.Logger.Debug("image removed", zap.String("ID", container.ImageID))
		}
	}

	if opts.RemoveNginx && containersInfo.NginxNetwork.ID != "" {
		err = cli.NetworkRemove(ctx, containersInfo.NginxNetwork.ID)
		if err != nil {
			return errcode.ErrDockerAPINetworkRemove.Wrap(err)
		}
		opts.Logger.Debug("network removed", zap.String("ID", containersInfo.NginxNetwork.ID))
	}

	return nil
}

func PS(ctx context.Context, depth int, cli *client.Client, logger *zap.Logger) error {
	logger.Debug("ps", zap.Int("depth", depth))

	containersInfo, err := GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NGINX ID", "STATUS", "CREATED"})
	if containersInfo.NginxContainer.ID != "" {
		table.Append([]string{
			containersInfo.NginxContainer.ID[:7],
			strings.Replace(containersInfo.NginxContainer.Status, "Up ", "", 1),
			strings.Replace(humanize.Time(time.Unix(containersInfo.NginxContainer.Created, 0)), " ago", "", 1),
		})
	}
	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "CHALLENGE", "SVC", "PORTS", "STATUS", "CREATED"})

	for _, flavor := range containersInfo.RunningFlavors {
		for uid, container := range flavor.Containers {
			ports := []string{}
			for _, port := range container.Ports {
				if port.PublicPort != 0 {
					ports = append(ports, strconv.Itoa(int(port.PublicPort)))
				}
			}

			table.Append([]string{
				uid[:7],
				flavor.ChallengeID(),
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
		RunningFlavors:    map[string]challengeFlavors{},
		RunningContainers: map[string]container{},
	}

	for _, dockerContainer := range containers {
		c := container(dockerContainer)

		// pathwar nginx proxy
		for _, name := range c.Names {
			if name[1:] == NginxContainerName {
				containersInfo.NginxContainer = c
			}
		}

		if _, found := c.Labels[challengeNameLabel]; !found { // not a pathwar container
			continue
		}

		flavor := c.ChallengeID()
		if _, found := containersInfo.RunningFlavors[flavor]; !found {
			challengeFlavor := challengeFlavors{
				Containers: map[string]container{},
			}
			challengeFlavor.Name = c.Labels[challengeNameLabel]
			challengeFlavor.Version = c.Labels[challengeVersionLabel]
			challengeFlavor.InstanceKey = c.Labels[InstanceKeyLabel]
			containersInfo.RunningFlavors[flavor] = challengeFlavor
		}
		containersInfo.RunningFlavors[flavor].Containers[c.ID] = c
		containersInfo.RunningContainers[c.ID] = c
	}

	// find proxy network
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, errcode.ErrDockerAPINetworkList.Wrap(err)
	}
	for _, networkResource := range networks {
		if networkResource.Name == ProxyNetworkName {
			containersInfo.NginxNetwork = networkResource
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
