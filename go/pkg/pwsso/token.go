package pwsso

import (
	"bytes"
	"encoding/json"
	"fmt"
	io "io"
	"net/http"
	time "time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (c *client) TokenWithClaims(bearer string) (*jwt.Token, jwt.MapClaims, error) {
	token, claims, err := TokenWithClaims(bearer, c.publicKey, c.opts.AllowUnsafe)
	if err != nil {
		c.logger.Warn("token with claims",
			zap.Error(err),
			zap.Any("pubkey", c.publicKey),
			zap.Bool("allow-unsafe", c.opts.AllowUnsafe),
			zap.String("bearer", bearer),
		)
	}
	return token, claims, err
}

func TokenWithClaims(bearer string, pubkey interface{}, allowUnsafe bool) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	// FIXME: add an option to automatically fetch the public key from
	// https://id.pathwar.land/auth/realms/Pathwar-Dev/protocol/openid-connect/certs
	// or
	// https://id.pathwar.land/auth/realms/Pathwar-Dev

	kf := func(token *jwt.Token) (interface{}, error) {
		return pubkey, nil
	}
	token, err := jwt.ParseWithClaims(bearer, claims, kf)
	if err != nil {
		if allowUnsafe {
			zap.L().Warn(
				"invalid token",
				zap.Error(err),
				zap.Bool("client-unsafe", true),
			)
			parser := new(jwt.Parser)
			token, _, err := parser.ParseUnverified(bearer, claims)
			if err != nil {
				return nil, nil, errcode.ErrSSOInvalidBearer.Wrap(err)
			}
			return token, claims, nil
		}
		e, ok := err.(*jwt.ValidationError)
		if !ok || (ok && e.Errors&jwt.ValidationErrorIssuedAt == 0) { // don't report error that token used before issued.
			return nil, nil, errcode.ErrSSOInvalidBearer.Wrap(err)
		}
	}
	return token, claims, nil
}

func TokenHasRole(token *jwt.Token, expectedRole string) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errcode.TODO // Unable to get claims out of JWT token
	}
	permissions, ok := claims["permissions"].([]interface{})
	if !ok {
		return errcode.TODO // Unable to get resource_access of JWT token claims
	}
	for _, permission := range permissions {
		if permission == expectedRole {
			return nil
		}
	}

	return errcode.TODO // Unable to get expected role from JWT token claims
}

func SubjectFromToken(token *jwt.Token) string {
	mc := token.Claims.(jwt.MapClaims)
	if v := mc["sub"]; v != nil {
		return v.(string)
	}
	return ""
}

func ClaimsFromToken(token *jwt.Token) *Claims {
	mc := token.Claims.(jwt.MapClaims)
	claims := &Claims{
		ActionToken: &ActionToken{},
	}

	// Subject & Issue at & Expiration
	if v := mc["sub"]; v != nil {
		claims.ActionToken.Sub = v.(string)
	}
	if v := mc["iat"]; v != nil {
		t := time.Unix(int64(v.(float64)), 0)
		claims.ActionToken.Iat = &t
	}
	if v := mc["exp"]; v != nil {
		t := time.Unix(int64(v.(float64)), 0)
		claims.ActionToken.Exp = &t
	}

	//FIXME: add more claims
	return claims
}

func GetUserInfoFromToken(token *jwt.Token, claims *Claims) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ProviderUserInfoURL, &bytes.Buffer{})
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var userinfo map[string]interface{}
	err = json.Unmarshal(body, &userinfo)
	fmt.Println("body is", string(body))
	if err != nil {
		return err
	}

	if v := userinfo["preferred_username"]; v != nil {
		claims.PreferredUsername = v.(string)
	} else if v := userinfo["nickname"]; v != nil {
		claims.PreferredUsername = v.(string)
	}
	if v := userinfo["email"]; v != nil {
		claims.Email = v.(string)
	}
	if v := userinfo["email_verified"]; v != nil {
		claims.EmailVerified = v.(bool)
	}
	return nil
}
