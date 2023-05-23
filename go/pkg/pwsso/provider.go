package pwsso

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

const (
	ProviderAuthURL     = "https://%s/authorize"
	ProviderTokenURL    = "https://%s/oauth/token"
	ProviderUserInfoURL = "https://%s/userinfo"
	ProviderRedirectURL = "https://html-tests.netlify.app/qs/"
	ProviderAudience    = "https://pathwar.land/"
)

func (c *client) Whoami(token string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(ProviderUserInfoURL, c.realm), &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userinfo map[string]interface{}
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		return nil, errcode.ErrSSOInvalidProviderResponse.Wrap(err)
	}

	return userinfo, err
}
