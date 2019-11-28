package pwsso

import (
	time "time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
)

func (c *client) TokenWithClaims(bearer string) (*jwt.Token, jwt.MapClaims, error) {
	return TokenWithClaims(bearer, c.publicKey, c.opts.AllowUnsafe)
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
		return nil, nil, errcode.ErrSSOInvalidBearer.Wrap(err)
	}
	return token, claims, nil
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

	// keycloak
	if v := mc["typ"]; v != nil {
		claims.ActionToken.Typ = v.(string)
	}
	if v := mc["sub"]; v != nil {
		claims.ActionToken.Sub = v.(string)
	}
	if v := mc["azp"]; v != nil {
		claims.ActionToken.Azp = v.(string)
	}
	if v := mc["iss"]; v != nil {
		claims.ActionToken.Iss = v.(string)
	}
	if v := mc["aud"]; v != nil {
		claims.ActionToken.Aud = v.(string)
	}
	if v := mc["asid"]; v != nil {
		claims.ActionToken.Asid = v.(string)
	}
	if v := mc["nonce"]; v != nil {
		claims.ActionToken.Nonce = v.(string)
	}
	if v := mc["session_state"]; v != nil {
		claims.ActionToken.SessionState = v.(string)
	}
	if v := mc["scope"]; v != nil {
		claims.ActionToken.Scope = v.(string)
	}
	if v := mc["jti"]; v != nil {
		claims.ActionToken.Jti = v.(string)
	}
	if v := mc["nbf"]; v != nil {
		claims.ActionToken.Nbf = float32(v.(float64))
	}
	if v := mc["iat"]; v != nil {
		t := time.Unix(int64(v.(float64)), 0)
		claims.ActionToken.Iat = &t
	}
	if v := mc["exp"]; v != nil {
		t := time.Unix(int64(v.(float64)), 0)
		claims.ActionToken.Exp = &t
	}
	if v := mc["auth_time"]; v != nil {
		t := time.Unix(int64(v.(float64)), 0)
		claims.ActionToken.AuthTime = &t
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

	// FIXME: add more infos
	return claims
}
