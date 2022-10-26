package pwapi

import (
	"context"
)

func (svc *service) SeasonList(ctx context.Context, in *SeasonList_Input) (*SeasonList_Output, error) {

	var ret SeasonList_Output
	svc.db.Find(&ret.Items)

	return &ret, nil
}
