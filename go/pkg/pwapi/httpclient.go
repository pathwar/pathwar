package pwapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-querystring/query"
	"pathwar.land/v2/go/pkg/errcode"
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

func (c HTTPClient) AgentListInstances(input *AgentListInstances_Input) (AgentListInstances_Output, error) {
	var result AgentListInstances_Output
	err := c.doGet("/agent/list-instances", input, &result)
	return result, err
}

func (c HTTPClient) AdminRedump(input *AdminRedump_Input) (AdminRedump_Output, error) {
	var result AdminRedump_Output
	err := c.doPost("/admin/redump", input, &result)
	return result, err
}

func (c HTTPClient) AdminPS(input *AdminPS_Input) (AdminPS_Output, error) {
	var result AdminPS_Output
	err := c.doGet("/admin/ps", input, &result)
	return result, err
}

func (c HTTPClient) AdminAddChallenge(input *AdminChallengeAdd_Input) (AdminChallengeAdd_Output, error) {
	var result AdminChallengeAdd_Output
	err := c.doPost("/admin/add-challenge", input, &result)
	return result, err
}

func (c HTTPClient) doPost(path string, input, output proto.Message) error {
	marshaler := jsonpb.Marshaler{}
	inputString, err := marshaler.MarshalToString(input)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	b := bytes.NewBuffer([]byte(inputString))

	url := c.baseAPI + path
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errcode.TODO.Wrap(fmt.Errorf("invalid status code (%d): %q", resp.StatusCode, string(body)))
	}

	err = jsonpb.Unmarshal(resp.Body, output)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}

func (c HTTPClient) doGet(path string, input, output proto.Message) error {
	qs, err := query.Values(input)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	url := c.baseAPI + path + "?" + qs.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errcode.TODO.Wrap(fmt.Errorf("invalid status code (%d): %q", resp.StatusCode, string(body)))
	}

	err = jsonpb.Unmarshal(resp.Body, output)
	if err != nil {
		return errcode.TODO.Wrap(err)
	}

	return nil
}
