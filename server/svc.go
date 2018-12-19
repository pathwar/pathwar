package server

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"pathwar.pw/entity"
)

type svc struct {
	jwtKey []byte
}

func (s *svc) Ping(_ context.Context, _ *Void) (*Void, error) {
	return &Void{}, nil
}

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

func (s *svc) Session(ctx context.Context, _ *Void) (*entity.Session, error) {
	sess, err := sessionFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get context session")
	}
	return &sess, nil
}
