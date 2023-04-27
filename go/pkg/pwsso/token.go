package pwsso

import (
	time "time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (c *client) TokenWithClaims(bearer string) (*jwt.Token, jwt.MapClaims, error) {
	token, claims, err := TokenWithClaims(bearer, c.publicKey, c.publicKey2, c.opts.AllowUnsafe)
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

func TokenWithClaims(bearer string, pubkey interface{}, pubkey2 interface{}, allowUnsafe bool) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	// FIXME: add an option to automatically fetch the public key from
	// https://id.pathwar.land/auth/realms/Pathwar-Dev/protocol/openid-connect/certs
	// or
	// https://id.pathwar.land/auth/realms/Pathwar-Dev

	kf := func(token *jwt.Token) (interface{}, error) {
		return pubkey, nil
	}

	kf2 := func(token *jwt.Token) (interface{}, error) {
		return pubkey2, nil
	}

	token, err := jwt.ParseWithClaims(bearer, claims, kf)
	if err != nil {
		token, err = jwt.ParseWithClaims(bearer, claims, kf2)
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
	}
	return token, claims, nil
}

func TokenHasRole(token *jwt.Token, expectedRole string) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errcode.TODO // Unable to get claims out of JWT token
	}
	resources, ok := claims["resource_access"]
	if !ok {
		return errcode.TODO // Unable to get resource_access of JWT token claims
	}
	clients, ok := resources.(map[string]interface{})
	if !ok {
		return errcode.TODO // Unable to get resources out of JWT token claims
	}
	for _, clientMap := range clients {
		client, ok := clientMap.(map[string]interface{})
		if !ok {
			continue
		}
		roles, ok := client["roles"].([]interface{})
		if !ok {
			continue
		}
		for _, role := range roles {
			if role == expectedRole {
				return nil
			}
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

	// pathwar specific
	if v := mc["preferred_username"]; v != nil {
		claims.PreferredUsername = v.(string)
	}
	if v := mc["email"]; v != nil {
		claims.Email = v.(string)
	}
	if v := mc["email_verified"]; v != nil {
		claims.EmailVerified = v.(bool)
	}
	if v := mc["given_name"]; v != nil {
		claims.GivenName = v.(string)
	}
	if v := mc["family_name"]; v != nil {
		claims.FamilyName = v.(string)
	}

	//FIXME: add more claims
	return claims
}
