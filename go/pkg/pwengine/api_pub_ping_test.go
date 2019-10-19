package pwengine

import (
	"context"
	"reflect"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_GetInfo(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := engine.GetInfo(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := &Info{
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
