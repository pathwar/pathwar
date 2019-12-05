package pwapi

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestSvc_ToolInfo(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := svc.ToolInfo(ctx, nil)
	checkErr(t, "", err)
	expected := &GetInfo_Output{
		Version: "dev",
		Commit:  "n/a",
		BuiltAt: "n/a",
		BuiltBy: "n/a",
	}
	expected.Uptime = status.Uptime // may vary
	testSameDeep(t, "", expected, status)
}
