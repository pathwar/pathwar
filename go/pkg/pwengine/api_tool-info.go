package pwengine

import (
	"context"
	"time"

	"pathwar.land/go/pkg/pwversion"
)

func (e *engine) ToolInfo(context.Context, *GetInfo_Input) (*GetInfo_Output, error) {
	return &GetInfo_Output{
		Version: pwversion.Version,
		Commit:  pwversion.Commit,
		BuiltAt: pwversion.Date,
		BuiltBy: pwversion.BuiltBy,
		Uptime:  int32(time.Now().Sub(e.startedAt).Seconds()),
	}, nil
}
