package pwagent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"go.uber.org/zap"
	"moul.io/godev"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwapi"
	"pathwar.land/pathwar/v2/go/pkg/pwversion"
)

func agentRegister(ctx context.Context, apiClient *pwapi.HTTPClient, opts Opts) error {
	metadata := map[string]interface{}{
		"delay":      opts.LoopDelay,
		"once":       opts.RunOnce,
		"num_cpus":   runtime.NumCPU(),
		"go_version": runtime.Version(),
		"uid":        os.Getuid(),
		"pid":        os.Getpid(),
	}
	metadataStr, err := json.Marshal(metadata)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	nginxPort, _ := strconv.Atoi(opts.HostPort)
	ret, err := apiClient.AgentRegister(ctx, &pwapi.AgentRegister_Input{
		Name:         opts.Name,
		Hostname:     hostname,
		NginxPort:    int32(nginxPort),
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		Version:      pwversion.Version,
		Tags:         []string{},
		DomainSuffix: opts.DomainSuffix,
		AuthSalt:     opts.AuthSalt,
		Metadata:     string(metadataStr),
		DefaultAgent: opts.DefaultAgent,
	})
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		fmt.Fprintln(os.Stderr, godev.PrettyJSON(ret))
	}
	return nil
}
