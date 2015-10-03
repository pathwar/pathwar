package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PathwarGenerateAToken struct {
	ID      string `json:"_id"`
	Created string `json:"_created"`
	Etag    string `json:"_etag"`
	Status  string `json:"_status"`
}

func (p *APIPathwar) GenerateAToken(login, password string, tmp bool) (*PathwarGenerateAToken, error) {
	request := p.client.Post(strings.Join([]string{APIUrl, "user-tokens"}, "/"))
	request = request.SetBasicAuth(login, password)

	request = request.Send(fmt.Sprintf("{\"is_session\": %v}", tmp))
	if p.debug {
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

func (p *APIPathwar) GetToken(login, password, id string) (*PathwarToken, error) {
	request := p.client.Get(fmt.Sprintf("%s/user-tokens/%s", APIUrl, id))
	request = request.SetBasicAuth(login, password)
	if p.debug {
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
