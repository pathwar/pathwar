package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/parnurzeal/gorequest"
)

func marshalWhere(where interface{}) (string, error) {
	if where == nil {
		return "{}", nil
	}

	whereString, err := json.Marshal(where)
	if err != nil {
		return "", err
	}
	return string(whereString), nil
}

func (p *APIPathwar) GetUsers(where interface{}) (*Users, error) {
	whereString, err := marshalWhere(where)
	if err != nil {
		return nil, err
	}

	resp, err := p.GetRequest(fmt.Sprintf("users?where=%s", whereString))
	var result Users
	err = json.Unmarshal(resp, &result)
	return &result, err
}

func (p *APIPathwar) GetRawOrganizationUsers(where interface{}) (*RawOrganizationUsers, error) {
	whereString, err := marshalWhere(where)
	if err != nil {
		return nil, err
	}

	resp, err := p.GetRequest(fmt.Sprintf("raw-organization-users?where=%s", whereString))
	var result RawOrganizationUsers
	err = json.Unmarshal(resp, &result)
	return &result, err
}

func (p *APIPathwar) GetRawLevelInstanceUsers(where interface{}) (*RawLevelInstanceUsers, error) {
	whereString, err := marshalWhere(where)
	if err != nil {
		return nil, err
	}

	resp, err := p.GetRequest(fmt.Sprintf("raw-level-instance-users?where=%s", whereString))
	var result RawLevelInstanceUsers
	err = json.Unmarshal(resp, &result)
	return &result, err
}

func (p *APIPathwar) GetRawLevelInstances(where interface{}) (*RawLevelInstances, error) {
	whereString, err := marshalWhere(where)
	if err != nil {
		return nil, err
	}

	resp, err := p.GetRequest(fmt.Sprintf("raw-level-instances?where=%s", whereString))
	var result RawLevelInstances
	err = json.Unmarshal(resp, &result)
	return &result, err
}

// ---

type PathwarGenerateAToken struct {
	ID      string `json:"_id"`
	Created string `json:"_created"`
	Etag    string `json:"_etag"`
	Status  string `json:"_status"`
}

func GenerateAToken(login, password string, tmp bool) (*PathwarGenerateAToken, error) {
	request := gorequest.New().Post(fmt.Sprintf("%s/user-tokens/", APIUrl))
	request = request.SetBasicAuth(login, password)
	request = request.Send(fmt.Sprintf("{\"is_session\": %v}", tmp))
	if os.Getenv("PATHWAR_DEBUG") != "" {
		request = request.SetDebug(true)
	}
	resp, body, errs := request.EndBytes()

	if len(errs) != 0 {
		return nil, printErrors(errs)
	}
	if err := httpHandleError([]int{201}, resp.StatusCode, body); err != nil {
		return nil, err
	}
	ret := &PathwarGenerateAToken{}

	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type PathwarToken struct {
	Token string `json:"token"`
}

func GetToken(login, password, id string) (*PathwarToken, error) {
	request := gorequest.New().Get(fmt.Sprintf("%s/user-tokens/%s", APIUrl, id))
	request = request.SetBasicAuth(login, password)
	if os.Getenv("PATHWAR_DEBUG") != "" {
		request = request.SetDebug(true)
	}
	resp, body, errs := request.EndBytes()

	if len(errs) != 0 {
		return nil, printErrors(errs)
	}
	if err := httpHandleError([]int{200}, resp.StatusCode, body); err != nil {
		return nil, err
	}
	ret := &PathwarToken{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
