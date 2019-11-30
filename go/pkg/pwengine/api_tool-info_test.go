package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_ToolInfo(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := engine.ToolInfo(ctx, nil)
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
