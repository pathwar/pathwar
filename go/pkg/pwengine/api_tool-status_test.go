package pwengine

import (
	"context"
	"testing"

	"pathwar.land/go/internal/testutil"
)

func TestEngine_GetStatus(t *testing.T) {
	engine, cleanup := TestingEngine(t, Opts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := engine.ToolStatus(ctx, nil)
	checkErr(t, "", err)

	expected := &GetStatus_Output{
		EverythingIsOK: true,
	}

	testSameDeep(t, "", expected, status)
}
