package pwapi

import (
	"bytes"
	"context"
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

func (c HTTPClient) AdminListChallenges(ctx context.Context, input *AdminListChallenges_Input) (AdminListChallenges_Output, error) {
	var _ *AdminListChallenges_Input = input
	var result AdminListChallenges_Output
	err := c.doGet(ctx, "/admin/list-challenges", input, &result)
	return result, err
}

func (c HTTPClient) AdminListUsers(ctx context.Context, input *AdminListUsers_Input) (AdminListUsers_Output, error) {
	var _ *AdminListUsers_Input = input
	var result AdminListUsers_Output
	err := c.doGet(ctx, "/admin/list-users", input, &result)
	return result, err
}

func (c HTTPClient) AdminListOrganizations(ctx context.Context, input *AdminListOrganizations_Input) (AdminListOrganizations_Output, error) {
	var _ *AdminListOrganizations_Input = input
	var result AdminListOrganizations_Output
	err := c.doGet(ctx, "/admin/list-organizations", input, &result)
	return result, err
}

func (c HTTPClient) AdminListTeams(ctx context.Context, input *AdminListTeams_Input) (AdminListTeams_Output, error) {
	var _ *AdminListTeams_Input = input
	var result AdminListTeams_Output
	err := c.doGet(ctx, "/admin/list-teams", input, &result)
	return result, err
}

func (c HTTPClient) AdminListAgents(ctx context.Context, input *AdminListAgents_Input) (AdminListAgents_Output, error) {
	var _ *AdminListAgents_Input = input
	var result AdminListAgents_Output
	err := c.doGet(ctx, "/admin/list-agents", input, &result)
	return result, err
}

func (c HTTPClient) AdminListCoupons(ctx context.Context, input *AdminListCoupons_Input) (AdminListCoupons_Output, error) {
	var _ *AdminListCoupons_Input = input
	var result AdminListCoupons_Output
	err := c.doGet(ctx, "/admin/list-coupons", input, &result)
	return result, err
}

func (c HTTPClient) AdminListChallengeSubscriptions(ctx context.Context, input *AdminListChallengeSubscriptions_Input) (AdminListChallengeSubscriptions_Output, error) {
	var _ *AdminListChallengeSubscriptions_Input = input
	var result AdminListChallengeSubscriptions_Output
	err := c.doGet(ctx, "/admin/list-challenge-subscriptions", input, &result)
	return result, err
}

func (c HTTPClient) AdminListActivities(ctx context.Context, input *AdminListActivities_Input) (AdminListActivities_Output, error) {
	var _ *AdminListActivities_Input = input
	var result AdminListActivities_Output
	err := c.doGet(ctx, "/admin/list-activities", input, &result)
	return result, err
}

func (c HTTPClient) AdminListAll(ctx context.Context, input *AdminListAll_Input) (AdminListAll_Output, error) {
	var _ *AdminListAll_Input = input
	var result AdminListAll_Output
	err := c.doGet(ctx, "/admin/list-all", input, &result)
	return result, err
}

func (c HTTPClient) AdminSearch(ctx context.Context, input *AdminSearch_Input) (AdminSearch_Output, error) {
	var _ *AdminSearch_Input = input
	var result AdminSearch_Output
	err := c.doPost(ctx, "/admin/search", input, &result)
	return result, err
}

func (c HTTPClient) AdminAddCoupon(ctx context.Context, input *AdminAddCoupon_Input) (AdminAddCoupon_Output, error) {
	var _ *AdminAddCoupon_Input = input
	var result AdminAddCoupon_Output
	err := c.doPost(ctx, "/admin/add-coupon", input, &result)
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

func (c HTTPClient) AdminAddSeasonChallenge(ctx context.Context, input *AdminSeasonChallengeAdd_Input) (AdminSeasonChallengeAdd_Output, error) {
	var _ *AdminSeasonChallengeAdd_Input = input
	var result AdminSeasonChallengeAdd_Output
	err := c.doPost(ctx, "/admin/season-challenge-add", input, &result)
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
	inputString := ""
	if input != nil {
		marshaler := jsonpb.Marshaler{}
		var err error
		inputString, err = marshaler.MarshalToString(input)
		if err != nil {
			return errcode.TODO.Wrap(err)
		}
	}

	ret, err := c.Raw(ctx, method, path, []byte(inputString))
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
