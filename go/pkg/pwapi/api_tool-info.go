package pwapi

import (
	"context"
	"time"

	"pathwar.land/go/pkg/pwversion"
)

func (svc *service) ToolInfo(context.Context, *GetInfo_Input) (*GetInfo_Output, error) {
	return &GetInfo_Output{
		Version: pwversion.Version,
		Commit:  pwversion.Commit,
		BuiltAt: pwversion.Date,
		BuiltBy: pwversion.BuiltBy,
		Uptime:  int32(time.Now().Sub(svc.startedAt).Seconds()),
	}, nil
}
