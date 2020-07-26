package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListUsers(ctx context.Context, in *AdminListUsers_Input) (*AdminListUsers_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var users []*pwdb.User
	err := svc.db.
		// Preload("").
		Find(&users).Error
	if err != nil {
		return nil, errcode.ErrListUsers.Wrap(err)
	}

	out := AdminListUsers_Output{Users: users}
	return &out, nil
}
