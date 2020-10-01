package testutil

import (
	"flag"
	"testing"

	"go.uber.org/zap"
	"moul.io/zapconfig"
)

var debug = flag.Bool("debug", false, "more verbose logging")

func Logger(t *testing.T) *zap.Logger {
	t.Helper()
	if !*debug {
		return zap.NewNop()
	}

	logger, err := zapconfig.Configurator{}.Build()
	if err != nil {
		t.Errorf("debug logger: %v", err)
		return zap.NewNop()
	}
	return logger
}
