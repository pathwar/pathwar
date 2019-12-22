package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/go/internal/testutil"
)

func TestSvc_GetStatus(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := context.Background()

	status, err := svc.ToolStatus(ctx, nil)
	checkErr(t, "", err)

	expected := &GetStatus_Output{
		EverythingIsOK: true,
	}

	assert.Equal(t, expected, status)
}
