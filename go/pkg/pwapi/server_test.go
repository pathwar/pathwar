package pwapi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"moul.io/godev"
	"pathwar.land/v2/go/internal/testutil"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := testutil.Logger(t)

	// init server
	server, cleanup := TestingServer(t, ctx, ServerOpts{Logger: logger})
	defer cleanup()

	{ // http
		svc := fmt.Sprintf("http://%s", server.ListenerAddr())
		resp, err := http.Get(svc + "/status")
		assert.NoError(t, err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		expected := "{\n  \"everything_is_ok\": true\n}"
		assert.Equal(t, expected, string(body))
		// FIXME: check rest of the headers (CORS, Content-Type, etc)
	}

	{ // gRPC
		client, cleanup := TestingClient(t, server.ListenerAddr())
		defer cleanup()
		ret, err := client.ToolStatus(ctx, &GetStatus_Input{})
		assert.NoError(t, err)
		assert.NotNil(t, ret, func() {
			fmt.Println(godev.PrettyJSON(ret))
			assert.Equal(t, true, ret.EverythingIsOK)
		})
	}
}
