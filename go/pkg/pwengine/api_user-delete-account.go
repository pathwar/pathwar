package pwengine

import (
	"context"
	"fmt"

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
		OAuthSubject:   fmt.Sprintf("deleted-%s", user.OAuthSubject),
		DeletionReason: in.Reason,
	}
	err = e.db.Model(&user).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// FIXME: invalide current JWT token

	ret := &UserDeleteAccountOutput{}
	return ret, nil
}
