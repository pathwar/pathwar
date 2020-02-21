package pwagent

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwapi"
)

func fetchAPIInstances(ctx context.Context, apiClient *http.Client, httpAPIAddr string, agentName string, logger *zap.Logger) (*pwapi.AgentListInstances_Output, error) {
	var instances pwapi.AgentListInstances_Output

	resp, err := apiClient.Get(httpAPIAddr + "/agent/list-instances?agent_name=" + agentName)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("received API error", zap.String("body", string(body)), zap.Int("code", resp.StatusCode))
		return nil, errcode.TODO.Wrap(fmt.Errorf("received API error"))
	}
	if err := jsonpb.UnmarshalString(string(body), &instances); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return &instances, nil
}
