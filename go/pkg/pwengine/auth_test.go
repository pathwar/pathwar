package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/pkg/pwsso"
)

func testSetContextToken(t *testing.T, ctx context.Context) context.Context {
	t.Helper()

	return context.WithValue(ctx, userTokenCtx, pwsso.TestingToken(t))
}
