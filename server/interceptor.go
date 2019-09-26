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

	"pathwar.land/client"
)

type ctxKey string

var (
	userTokenCtx ctxKey = "user-token"
)

func (s *svc) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	switch fullMethodName {
	// do not check for token for the following services
	case "/pathwar.server.Server/Authenticate", "/pathwar.server.Server/Ping", "/pathwar.server.Server/Info", "/pathwar.server.Server/Status":
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

	token, _, err := client.TokenWithClaims(auth[0], s.client)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	ctx = context.WithValue(ctx, userTokenCtx, token)

	return ctx, nil
}

func authFunc(ctx context.Context) (context.Context, error) {
	// do nothing here, use AuthFuncOverride instead
	return nil, fmt.Errorf("should not happen")
}

// func claimsFromContext(ctx context.Context) (*client.Claims, error) {}

func subjectFromContext(ctx context.Context) (string, error) {
	token, err := userTokenFromContext(ctx)
	if err != nil {
		return "", err
	}

	sub := client.SubjectFromToken(token)
	if sub == "" {
		return "", errors.New("no such subject in the token")
	}

	return sub, nil
}

func userTokenFromContext(ctx context.Context) (*jwt.Token, error) {
	token := ctx.Value(userTokenCtx)
	if token == nil {
		return nil, errors.New("context does not contain a session")
	}
	return token.(*jwt.Token), nil
}
