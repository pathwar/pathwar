package pwengine

import (
	"context"
	"reflect"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_ToolInfo(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := engine.ToolInfo(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := &GetInfo_Output{
		Version: "dev",
		Commit:  "n/a",
		BuiltAt: "n/a",
		BuiltBy: "n/a",
	}
	expected.Uptime = status.Uptime // may vary
	if !reflect.DeepEqual(expected, status) {
		t.Fatalf("Expected: %#v, got %#v instead.", expected, status)
	}
}
