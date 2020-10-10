package pwsso

import (
	"fmt"
	"net/url"

	"github.com/keycloak/kcinit/rest"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

const KeycloakBaseURL = "https://id.pathwar.land"

func (c *client) Whoami(token string) (map[string]interface{}, error) {
	oidc, err := c.oidc()
	if err != nil {
		return nil, errcode.ErrSSOGetOIDC.Wrap(err)
	}

	res, err := oidc.Path("userinfo").Request().Header("Authorization", "bearer "+token).Get()
	if err != nil {
		return nil, errcode.ErrSSOFailedKeycloakRequest.Wrap(err)
	}

	var info map[string]interface{}
	if err := res.ReadJson(&info); err != nil {
		return nil, errcode.ErrSSOInvalidKeycloakResponse.Wrap(err)
	}

	return info, nil
}

func (c *client) Logout(token string) error {
	oidc, err := c.oidc()
	if err != nil {
		return errcode.ErrSSOGetOIDC.Wrap(err)
	}

	form := url.Values{}
	form.Set("client_id", c.clientID)
	// form.Set("client_secret", c.clientSecret)
	form.Set("refresh_token", token)
	res, err := oidc.Path("logout").Request().Form(form).Post()
	if err != nil {
		return errcode.ErrSSOLogout.Wrap(err)
	}
	var ret map[string]interface{}
	if err := res.ReadJson(&ret); err != nil {
		return errcode.ErrSSOInvalidKeycloakResponse.Wrap(err)
	}
	c.logger.Debug("keycloak returned", zap.Any("ret", ret))
	if _, ok := ret["error"]; ok {
		return errcode.ErrSSOKeycloakError.Wrap(fmt.Errorf("%s: %s", ret["error"].(string), ret["error_description"].(string)))
	}
	return nil
}

func (c *client) oidc() (*rest.WebTarget, error) {
	keycloak := rest.New()
	realmURL := fmt.Sprintf("%s/auth/realms/%s", KeycloakBaseURL, c.realm)
	base := keycloak.Target(realmURL)
	if base == nil {
		return nil, errcode.ErrSSOInitKeycloakClient
	}
	oidc := base.Path("protocol/openid-connect")
	return oidc, nil
}
