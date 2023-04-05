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

	// retrieve number of challenges solved by each team
	out := AdminSeasonStats_Output{}
	for i, team := range teams {
		stat := AdminSeasonStats_Output_Stat{Team: &teams[i], ChallengesSolved: 0}
		err = svc.db.
			Model(&pwdb.ChallengeValidation{}).
			Where(&pwdb.ChallengeValidation{TeamID: team.ID}).
			Count(&stat.ChallengesSolved).
			Error
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}
		out.Stats = append(out.Stats, &stat)
	}

	return &out, nil
}
