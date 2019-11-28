package pwengine

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/metadata"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

type ctxKey string

var (
	userTokenCtx ctxKey = "user-token"
)

// AuthFuncOverride exists to implement the grpc_auth.ServiceAuthFuncOverride interface
//
// see https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/auth/auth.go
func (e *engine) AuthFuncOverride(ctx context.Context, path string) (context.Context, error) {
	switch path { // always allow public endpoints
	case "/pathwar.engine.Engine/ToolPing",
		"/pathwar.engine.Engine/ToolStatus",
		"/pathwar.engine.Engine/ToolInfo":
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errcode.ErrMissingContextMetadata
	}

	auth, ok := md["authorization"]
	if !ok || len(auth) < 1 {
		return nil, errcode.ErrNoTokenProvided
	}

	// cleanup the authorization
	if strings.HasPrefix(auth[0], "Bearer ") {
		auth[0] = auth[0][7:]
	}

	token, _, err := e.sso.TokenWithClaims(auth[0])
	if err != nil {
		return nil, errcode.ErrGetTokenWithClaims.Wrap(err)
	}
	ctx = context.WithValue(ctx, userTokenCtx, token)
	return ctx, nil
}

func tokenFromContext(ctx context.Context) (*jwt.Token, error) {
	token := ctx.Value(userTokenCtx)
	if token == nil {
		return nil, errcode.ErrNoTokenInContext
	}
	return token.(*jwt.Token), nil
}

func subjectFromContext(ctx context.Context) (string, error) {
	token, err := tokenFromContext(ctx)
	if err != nil {
		return "", errcode.ErrGetTokenFromContext.Wrap(err)
	}

	sub := pwsso.SubjectFromToken(token)
	if sub == "" {
		return "", errcode.ErrGetSubjectFromToken
	}

	return sub, nil
}

func userIDFromContext(ctx context.Context, db *gorm.DB) (int64, error) {
	oauthSubject, err := subjectFromContext(ctx)
	if err != nil {
		return 0, errcode.ErrGetSubjectFromContext.Wrap(err)
	}

	// FIXME: only fetch the ID instead of the whole user
	var user pwdb.User
	err = db.
		Where(pwdb.User{OAuthSubject: oauthSubject}).
		Find(&user).
		Error
	if err != nil {
		return 0, errcode.ErrGetUserBySubject.Wrap(err)
	}

	return user.ID, nil
}
