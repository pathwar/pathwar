package pwengine

import (
	"context"
	"reflect"
	"testing"
)

func TestEngine_GetStatus(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{})
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
