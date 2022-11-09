package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) SeasonList(ctx context.Context, in *SeasonList_Input) (*SeasonList_Output, error) {
	var (
		ret SeasonList_Output
	)

	seasons, err := svc.loadUserSeasons(ctx)
	if err != nil {
		return nil, errcode.ErrLoadUserSeasons
	}

	for _, season := range seasons {
		ret.Seasons = append(ret.Seasons, &SeasonList_Output_SeasonAndTeam{season.Season, season.Team, season.IsActive})
	}

	return &ret, nil
}
