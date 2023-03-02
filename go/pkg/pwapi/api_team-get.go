package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) TeamGet(ctx context.Context, in *TeamGet_Input) (*TeamGet_Output, error) {
	if in == nil || in.TeamID == 0 {
		return nil, errcode.ErrMissingInput
	}

	var item pwdb.Team
	err := svc.db.
		Preload("Season").
		Preload("Organization").
		Preload("Members").
		Preload("Members.User"). // FIXME: only if member of the team or if admin
		Preload("ChallengeSubscriptions"). // FIXME: only if member of the team or if admin
		Preload("Achievements").
		Where(pwdb.Team{
			ID:             in.TeamID,
			DeletionStatus: pwdb.DeletionStatus_Active,
		}).
		First(&item).
		Error
	if err != nil {
		return nil, errcode.ErrGetTeam.Wrap(err)
	}

	ret := TeamGet_Output{Item: &item}

	return &ret, nil
}
