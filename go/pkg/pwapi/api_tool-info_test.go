package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/go/v2/internal/testutil"
)

func TestEngine_ToolInfo(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	status, err := svc.ToolInfo(ctx, nil)
	assert.NoError(t, err)
	expected := &GetInfo_Output{
		Version: "dev",
		Commit:  "n/a",
		BuiltAt: "n/a",
		BuiltBy: "n/a",
	}
	expected.Uptime = status.Uptime // may vary
	assert.Equal(t, expected, status)
}
