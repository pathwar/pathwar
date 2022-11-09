package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) SeasonList(ctx context.Context, in *SeasonList_Input) (*SeasonList_Output, error) {
	var (
		ret SeasonList_Output
		err error
	)

	ret.Seasons, err = svc.loadListSeasons(ctx)
	if err != nil {
		return nil, errcode.ErrLoadUserSeasons
	}

	return &ret, nil
}

func (svc *service) loadListSeasons(ctx context.Context) ([]*SeasonList_Output_SeasonAndTeam, error) {
	var (
		seasons     []*pwdb.Season
		memberships []*pwdb.TeamMember
	)

	userID, err := userIDFromContext(ctx, svc.db)
	if err != nil {
		return nil, errcode.ErrUnauthenticated
	}

	// get season organizations for user
	err = svc.db.
		Preload("Team").
		Preload("Team.Organization").
		Where(pwdb.TeamMember{UserID: userID}).
		Find(&memberships).
		Error
	if err != nil {
		return nil, errcode.ErrGetUserOrganizations.Wrap(err)
	}

	req := svc.db
	switch {
	case isAdminContext(ctx): // because it's the highest level
		// noop
	case isTesterContext(ctx):
		req = req.
			Where(pwdb.Season{Visibility: pwdb.Season_Public}).
			Or(pwdb.Season{IsTesting: true})
	default: // "normal" user
		req = req.
			Where(pwdb.Season{Visibility: pwdb.Season_Public})
	}

	// get all available seasons
	err = req.
		Find(&seasons).
		Error
	if err != nil {
		return nil, errcode.ErrGetSeasons.Wrap(err)
	}

	var output []*SeasonList_Output_SeasonAndTeam
	for _, season := range seasons {
		item := &SeasonList_Output_SeasonAndTeam{
			Season: season,
		}

		for _, membership := range memberships {
			if membership.Team.SeasonID == season.ID {
				item.Team = membership.Team
				break
			}
		}

		output = append(output, item)
	}

	return output, nil
}
