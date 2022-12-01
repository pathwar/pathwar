package pwes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestCompute(t *testing.T) {

	tests := []struct {
		name          string
		expectedErr   error
		expectedScore int
	}{
		{"api null", errcode.ErrMissingInput, 0},
	}
	var ctx context.Context
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := RollbackScore(ctx, nil)
			assert.Equal(t, test.expectedErr, err, test.name)
		})
	}
}
