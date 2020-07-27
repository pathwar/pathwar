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
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"pathwar.land/pathwar/v2/go/internal/randstring"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwinit"
)

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
		err = Clean(ctx, cli, CleanOpts{Logger: opts.Logger, ContainerIDs: []string{challengeID}})
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
