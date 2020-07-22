package pwapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-querystring/query"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func NewHTTPClient(httpClient *http.Client, baseAPI string) *HTTPClient {
	return &HTTPClient{
		baseAPI:    baseAPI,
		httpClient: httpClient,
	}
}

type HTTPClient struct {
	baseAPI    string
	httpClient *http.Client
}

func (c HTTPClient) AgentListInstances(ctx context.Context, input *AgentListInstances_Input) (AgentListInstances_Output, error) {
	var _ *AgentListInstances_Input = input
	var result AgentListInstances_Output
	err := c.doGet(ctx, "/agent/list-instances", input, &result)
	return result, err
}

func (c HTTPClient) AgentRegister(ctx context.Context, input *AgentRegister_Input) (AgentRegister_Output, error) {
	var _ *AgentRegister_Input = input
	var result AgentRegister_Output
	err := c.doPost(ctx, "/agent/register", input, &result)
	return result, err
}

func (c HTTPClient) AgentUpdateState(ctx context.Context, input *AgentUpdateState_Input) (AgentUpdateState_Output, error) {
	var _ *AgentUpdateState_Input = input
	var result AgentUpdateState_Output
	err := c.doPost(ctx, "/agent/update-state", input, &result)
	return result, err
}

func (c HTTPClient) AdminRedump(ctx context.Context, input *AdminRedump_Input) (AdminRedump_Output, error) {
	var _ *AdminRedump_Input = input
	var result AdminRedump_Output
	err := c.doPost(ctx, "/admin/redump", input, &result)
	return result, err
}

func (c HTTPClient) AdminPS(ctx context.Context, input *AdminPS_Input) (AdminPS_Output, error) {
	var _ *AdminPS_Input = input
	var result AdminPS_Output
	err := c.doGet(ctx, "/admin/ps", input, &result)
	return result, err
}

func (c HTTPClient) AdminAddChallenge(ctx context.Context, input *AdminChallengeAdd_Input) (AdminChallengeAdd_Output, error) {
	var _ *AdminChallengeAdd_Input = input
	var result AdminChallengeAdd_Output
	err := c.doPost(ctx, "/admin/challenge-add", input, &result)
	return result, err
}

func (c HTTPClient) AdminAddChallengeFlavor(ctx context.Context, input *AdminChallengeFlavorAdd_Input) (AdminChallengeFlavorAdd_Output, error) {
	var _ *AdminChallengeFlavorAdd_Input = input
	var result AdminChallengeFlavorAdd_Output
	err := c.doPost(ctx, "/admin/challenge-flavor-add", input, &result)
	return result, err
}

func (c HTTPClient) AdminAddChallengeInstance(ctx context.Context, input *AdminChallengeInstanceAdd_Input) (AdminChallengeInstanceAdd_Output, error) {
	var _ *AdminChallengeInstanceAdd_Input = input
	var result AdminChallengeInstanceAdd_Output
	err := c.doPost(ctx, "/admin/challenge-instance-add", input, &result)
	return result, err
}

func (c HTTPClient) GetStatus(ctx context.Context, input *GetStatus_Input) (GetStatus_Output, error) {
	var _ *GetStatus_Input = input
	var result GetStatus_Output
	err := c.doGet(ctx, "/status", input, &result)
	return result, err
}

func (c HTTPClient) UserSetPreferences(ctx context.Context, input *UserSetPreferences_Input) (UserSetPreferences_Output, error) {
	var _ *UserSetPreferences_Input = input
	var result UserSetPreferences_Output
	err := c.doPost(ctx, "/user/preferences", input, &result)
	return result, err
}

func (c HTTPClient) RawProto(ctx context.Context, method string, path string, input, output proto.Message) error {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	ret, err := c.Raw(ctx, method, path, inputBytes)
	if err != nil {
		return err
	}

	return jsonpb.UnmarshalString(string(ret), output)
}

func (c HTTPClient) Raw(ctx context.Context, method string, path string, input []byte) ([]byte, error) {
	url := c.baseAPI + path
	b := bytes.NewBuffer(input)

	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errcode.TODO.Wrap(fmt.Errorf("invalid status code (%d): %q", resp.StatusCode, string(body)))
	}

	return ioutil.ReadAll(resp.Body)
}

func (c HTTPClient) doPost(ctx context.Context, path string, input, output proto.Message) error {
	marshaler := jsonpb.Marshaler{}
	inputString, err := marshaler.MarshalToString(input)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	ret, err := c.Raw(ctx, "POST", path, []byte(inputString))
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	b := bytes.NewBuffer(ret)
	err = jsonpb.Unmarshal(b, output)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}

func (c HTTPClient) doGet(ctx context.Context, path string, input, output proto.Message) error {
	qs, err := query.Values(input)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	path = path + "?" + qs.Encode()

	ret, err := c.Raw(ctx, "GET", path, nil)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	b := bytes.NewBuffer(ret)
	err = jsonpb.Unmarshal(b, output)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}
