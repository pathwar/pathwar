package pwapi

import (
	"context"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"testing"
)

func TestService_SeasonList(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// fetch user session to ensure account is created
	_, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
}
