package pwsso

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/keycloak/kcinit/rest"
	"go.uber.org/zap"
)

const (
	keycloakBaseURL = "https://id.pathwar.land/auth"
)

func (c *client) Whoami(token string) (map[string]interface{}, error) {
	oidc, err := c.oidc()
	if err != nil {
		return nil, fmt.Errorf("failed to get oidc: %w", err)
	}

	res, err := oidc.Path("userinfo").Request().Header("Authorization", "brear "+token).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo from keycloak: %w", err)
	}

	var info map[string]interface{}
	if err := res.ReadJson(&info); err != nil {
		return nil, fmt.Errorf("failed to read JSON from keycloak response: %w", err)
	}

	return info, nil
}

func (c *client) Logout(token string) error {
	oidc, err := c.oidc()
	if err != nil {
		return fmt.Errorf("failed to get oidc: %w", err)
	}

	form := url.Values{}
	form.Set("client_id", c.clientID)
	//form.Set("client_secret", c.clientSecret)
	form.Set("refresh_token", token)
	res, err := oidc.Path("logout").Request().Form(form).Post()
	if err != nil {
		return fmt.Errorf("failed to logout from keycloak: %w", err)
	}
	var ret map[string]interface{}
	if err := res.ReadJson(&ret); err != nil {
		return fmt.Errorf("failed to read result from keycloak: %w", err)
	}
	c.logger.Debug("keycloak returned", zap.Any("ret", ret))
	if _, ok := ret["error"]; ok {
		return fmt.Errorf("%s: %s", ret["error"].(string), ret["error_description"].(string))
	}
	return nil
}

func (c *client) oidc() (*rest.WebTarget, error) {
	keycloak := rest.New()
	realmURL := fmt.Sprintf("%s/realms/%s", keycloakBaseURL, c.realm)
	base := keycloak.Target(realmURL)
	if base == nil {
		return nil, errors.New("failed to initialize keycloak client")
	}
	oidc := base.Path("protocol/openid-connect")
	return oidc, nil
}
