package pwapi

import (
	"context"
	"strconv"
	"strings"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminSeasonStats(ctx context.Context, in *AdminSeasonStats_Input) (*AdminSeasonStats_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil || in.SeasonID == "" {
		return nil, errcode.ErrMissingInput
	}

	seasonID, err := pwdb.GetIDBySlugAndKind(svc.db, in.SeasonID, "season")
	if err != nil {
		return nil, err
	}

	teams := []pwdb.Team{}
	err = svc.db.
		Preload("Members").
		Preload("Members.User").
		Where(&pwdb.Team{SeasonID: seasonID}).
		Where("score > 0").
		Order("score DESC").
		Find(&teams).
		Error
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	// retrieve number of challenges solved by each team
	out := AdminSeasonStats_Output{}
	var challengesSolved int64
	for rank, team := range teams {
		err = svc.db.
			Model(&pwdb.ChallengeValidation{}).
			Where(&pwdb.ChallengeValidation{TeamID: team.ID}).
			Count(&challengesSolved).
			Error
		if err != nil {
			return nil, errcode.TODO.Wrap(err)
		}
		for _, member := range team.Members {
			// teamName keep only part before @
			stat := AdminSeasonStats_Output_Stat{
				Rank:             strconv.FormatInt(int64(rank+1), 10),
				Mail:             member.User.Email,
				Name:             member.User.Slug,
				TeamName:         team.Slug[:strings.LastIndex(team.Slug, "@")],
				Score:            strconv.FormatInt(team.Score, 10),
				ChallengesSolved: strconv.FormatInt(challengesSolved, 10),
			}
			out.Stats = append(out.Stats, &stat)
		}
	}

	return &out, nil
}
