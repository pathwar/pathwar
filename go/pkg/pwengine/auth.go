package pwengine

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"pathwar.land/go/pkg/pwsso"
)

type ctxKey string

var (
	userTokenCtx ctxKey = "user-token"
)

// AuthFuncOverride implements the grpc_auth.ServiceAuthFuncOverride interface
//
// see https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/auth/auth.go
func (c *client) AuthFuncOverride(ctx context.Context, path string) (context.Context, error) {
	switch path { // always allow public endpoints
	case "/pathwar.engine.Engine/Ping",
		"/pathwar.engine.Engine/GetStatus",
		"/pathwar.engine.Engine/GetInfo":
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "cannot get metadata from context")
	}

	auth, ok := md["authorization"]
	if !ok || len(auth) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "no token provided")
	}

	// cleanup the authorization
	if strings.HasPrefix(auth[0], "Bearer ") {
		auth[0] = auth[0][7:]
	}

	token, _, err := c.sso.TokenWithClaims(auth[0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	ctx = context.WithValue(ctx, userTokenCtx, token)

	return ctx, nil
}

func tokenFromContext(ctx context.Context) (*jwt.Token, error) {
	token := ctx.Value(userTokenCtx)
	if token == nil {
		return nil, errors.New("context does not contain a token")
	}
	return token.(*jwt.Token), nil
}

func subjectFromContext(ctx context.Context) (string, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token from contact: %w", err)
	}

	sub := pwsso.SubjectFromToken(token)
	if sub == "" {
		return "", errors.New("no such subject in the token")
	}

	return sub, nil
}
