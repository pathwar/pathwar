package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"pathwar.pw/entity"
)

type ctxKey string

var userSessionCtx ctxKey = "user-session"

func (s *svc) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	switch fullMethodName {
	// do not check for token for the following services
	case "/pathwar.server.Server/Authenticate", "/pathwar.server.Server/Ping":
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "cannot get metadata/headers")
	}
	auth, ok := md["authorization"]
	if !ok || len(auth) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "no token provided")
	}
	if strings.HasPrefix(auth[0], "Bearer ") {
		auth[0] = auth[0][7:]
	}

	token, err := jwt.Parse(auth[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "there was an error")
		}
		return s.jwtKey, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	ctx = context.WithValue(ctx, userSessionCtx, entity.UserSession{
		// FIXME: use mapstructure
		Username: claims["username"].(string),
	})

	return ctx, nil
}

func authFunc(ctx context.Context) (context.Context, error) {
	// do nothing here, use AuthFuncOverride instead
	return nil, fmt.Errorf("should not happen")
}

func userSessionFromContext(ctx context.Context) (entity.UserSession, error) {
	sess := ctx.Value(userSessionCtx)
	if sess == nil {
		return entity.UserSession{}, errors.New("context does not contain a session")
	}

	return sess.(entity.UserSession), nil
}
