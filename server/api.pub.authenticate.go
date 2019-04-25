package server

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
)

func (s *svc) Authenticate(ctx context.Context, input *AuthenticateInput) (*AuthenticateOutput, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// FIXME: use mapstructure
		"username": input.Username,
		// FIXME: if needed encrypt sensitive data
	})
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}
	// FIXME: set "cookie"
	return &AuthenticateOutput{
		Token: tokenString,
	}, nil
}
