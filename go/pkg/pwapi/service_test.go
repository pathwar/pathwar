package pwapi

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestService_impl(t *testing.T) {
	var _ Service = (*service)(nil)
	var _ ServiceServer = (*service)(nil)
}
