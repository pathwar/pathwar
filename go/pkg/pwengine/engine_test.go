package pwengine

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestEngine_impl(t *testing.T) {
	var _ Engine = (*engine)(nil)
	var _ EngineServer = (*engine)(nil)
}
