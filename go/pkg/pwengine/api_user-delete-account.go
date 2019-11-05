package pwengine

import (
	"context"
	"fmt"
	"time"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) UserDeleteAccount(ctx context.Context, in *UserDeleteAccountInput) (*UserDeleteAccountOutput, error) {
	userID, err := userIDFromContext(ctx, e.db)
	if err != nil {
		return nil, fmt.Errorf("get userid from context: %w", err)
	}

	var user pwdb.User
	err = e.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	updates := pwdb.User{
		OAuthSubject:   fmt.Sprintf("deleted_%s_%d", user.OAuthSubject, time.Now().Unix()),
		DeletionReason: in.Reason,
	}
	err = e.db.Model(&user).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// FIXME: mark the user state as deleted
	// FIXME: mark the solo team as deleted
	// FIXME: invalide current JWT token
	// FIXME: add another task that pseudonymize the data

	ret := &UserDeleteAccountOutput{}
	return ret, nil
}
