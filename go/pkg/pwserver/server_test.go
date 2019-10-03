package pwserver

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestServer(t *testing.T) {
	opts := Opts{
		HTTPBind: ":8000",
	}
	ctx, cancel := context.WithCancel(context.Background())
	start, cleanup := TestingServer(t, ctx, opts)
	defer cleanup()
	defer cancel()
	go start()

	resp, err := http.Get("http://localhost:8000/status")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := `{
  "everything_is_ok": true
}`
	if string(body) != expected {
		t.Fatalf("Expected %q, got %q instead.", expected, string(body))
	}
	// FIXME: check rest of the headers (CORS, Content-Type, etc)
}
