package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminSeasonStats(ctx context.Context, in *AdminSeasonStats_Input) (*AdminSeasonStats_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	seasonID, err := pwdb.GetIDBySlugAndKind(svc.db, in.SeasonID, "season")
	if err != nil {
		return nil, err
	}

	// retrieve all teams for the given season and preload team_member and team_member.user
	teams := []pwdb.Team{}
	err = svc.db.
		Preload("TeamMembers").
		Preload("TeamMembers.User").
		Where(&pwdb.Team{SeasonID: seasonID}).
		Find(&teams).
		Error
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return &AdminSeasonStats_Output{}, nil
}
