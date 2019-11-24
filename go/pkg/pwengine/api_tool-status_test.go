package pwengine

import (
	"context"
	"reflect"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_GetStatus(t *testing.T) {
	t.Parallel()
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := engine.ToolStatus(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := &GetStatus_Output{
		EverythingIsOK: true,
	}
	if !reflect.DeepEqual(expected, status) {
		t.Fatalf("Expected: %#v, got %#v instead.", expected, status)
	}
}
