package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_UserDeleteAccount(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// ensure account is created
	beforeDelete, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	beforeDeleteID := beforeDelete.User.ID
	beforeDeleteSubject := beforeDelete.User.OAuthSubject

	// delete account
	_, err = engine.UserDeleteAccount(ctx, &UserDeleteAccount_Input{Reason: "just a test"})
	if err != nil {
		t.Errorf("err: %v", err)
	}

	// create new account
	afterDelete, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if !afterDelete.IsNewUser {
		t.Errorf("Expected session.IsNewUser==true, got false.")
	}
	if beforeDeleteID == afterDelete.User.ID {
		t.Errorf("Expected different user id, got same.")
	}
	if beforeDeleteSubject != afterDelete.User.OAuthSubject {
		t.Errorf("Expected same OAuth subject, got different.")
	}

	// retrieve already created account
	afterAfterDelete, err := engine.UserGetSession(ctx, nil)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if afterAfterDelete.IsNewUser {
		t.Errorf("Expected session.IsNewUser==false, got true.")
	}
}
