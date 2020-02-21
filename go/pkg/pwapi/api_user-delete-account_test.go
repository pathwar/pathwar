package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/go/internal/testutil"
)

func TestSvc_UserDeleteAccount(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// ensure account is created
	beforeDelete, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
	beforeDeleteID := beforeDelete.User.ID
	beforeDeleteSubject := beforeDelete.User.OAuthSubject

	// delete account
	_, err = svc.UserDeleteAccount(ctx, &UserDeleteAccount_Input{Reason: "just a test"})
	require.NoError(t, err)

	// create new account
	afterDelete, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
	assert.True(t, afterDelete.IsNewUser)
	assert.NotEqual(t, beforeDeleteID, afterDelete.User.ID)
	assert.Equal(t, beforeDeleteSubject, afterDelete.User.OAuthSubject)

	// retrieve already created account
	afterAfterDelete, err := svc.UserGetSession(ctx, nil)
	require.NoError(t, err)
	assert.False(t, afterAfterDelete.IsNewUser)
}
