package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_UserDeleteAccount(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	// ensure account is created
	beforeDelete, err := engine.UserGetSession(ctx, nil)
	checkErr(t, "", err)
	beforeDeleteID := beforeDelete.User.ID
	beforeDeleteSubject := beforeDelete.User.OAuthSubject

	// delete account
	_, err = engine.UserDeleteAccount(ctx, &UserDeleteAccount_Input{Reason: "just a test"})
	checkErr(t, "", err)

	// create new account
	afterDelete, err := engine.UserGetSession(ctx, nil)
	checkErr(t, "", err)
	if !afterDelete.IsNewUser {
		t.Errorf("Expected session.IsNewUser==true, got false.")
	}
	testDifferentInt64s(t, "", beforeDeleteID, afterDelete.User.ID)
	testSameStrings(t, "", beforeDeleteSubject, afterDelete.User.OAuthSubject)

	// retrieve already created account
	afterAfterDelete, err := engine.UserGetSession(ctx, nil)
	checkErr(t, "", err)
	if afterAfterDelete.IsNewUser {
		t.Errorf("Expected session.IsNewUser==false, got true.")
	}
}
