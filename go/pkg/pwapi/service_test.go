package pwapi

import (
	"testing"
)

func TestService_impl(t *testing.T) {
	var _ Service = (*service)(nil)
	var _ ServiceServer = (*service)(nil)
}
