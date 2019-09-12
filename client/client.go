package client

import (
	time "time"

	jwt "github.com/dgrijalva/jwt-go"
)

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
