package pwengine

import (
	"context"
	"reflect"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_GetStatus(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := engine.GetStatus(ctx, nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expected := &Status{
		EverythingIsOK: true,
	}
	if !reflect.DeepEqual(expected, status) {
		t.Fatalf("Expected: %#v, got %#v instead.", expected, status)
	}
}
