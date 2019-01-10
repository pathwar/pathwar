package server

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"pathwar.pw/entity"
	"pathwar.pw/sql"
)

type svc struct {
	jwtKey []byte
	db     *gorm.DB
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

func (s *svc) GenerateFakeData(ctx context.Context, _ *Void) (*Void, error) {
	if err := s.db.Create(&entity.Level{
		Name:        "level1",
		Description: "description 1",
		Author:      "author 1",
		SourceUrl:   "http://1.2.3.4",
		Version:     "1.0",
		Lang:        "en",
		//Tags:        []string{"test", "example"},
		IsDraft: false,
	}).Error; err != nil {
		return nil, err
	}
	return &Void{}, nil
}

func (s *svc) Dump(ctx context.Context, _ *Void) (*entity.Dump, error) {
	return sql.DoDump(s.db)
}
