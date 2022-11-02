package pwapi

import (
	"context"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) SeasonList(ctx context.Context, in *SeasonList_Input) (*SeasonList_Output, error) {
	var ret SeasonList_Output
	svc.db.Where(pwdb.Season{IsGlobal: true, IsTesting: false, Status: 1}).Find(&ret.Items)

	return &ret, nil
}
