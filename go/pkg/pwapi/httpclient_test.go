package pwapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/pathwar/v2/go/internal/testutil"
	"pathwar.land/pathwar/v2/go/pkg/pwsso"
)

func TestHTTPClient_GetStatus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := testutil.Logger(t)

	// init server
	server, cleanup := TestingServer(t, ctx, ServerOpts{Logger: logger})
	defer cleanup()

	// init client
	bind := fmt.Sprintf("http://%s", server.ListenerAddr())
	hc := &http.Client{
		Transport: pwsso.TestingTransport(t),
	}
	client := NewHTTPClient(hc, bind)

	status, err := client.GetStatus(ctx, &GetStatus_Input{})
	require.NoError(t, err)
	expected := GetStatus_Output{EverythingIsOK: true}
	assert.Equal(t, expected, status)

	ret, err := client.UserSetPreferences(ctx, &UserSetPreferences_Input{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), `TODO(#666): TODO(#666): invalid status code (500): "{\n  \"code\": 2,\n  \"message\": \"ErrUnauthenticated(#103)\"\n}"`)
	assert.Equal(t, UserSetPreferences_Output{}, ret)
}
