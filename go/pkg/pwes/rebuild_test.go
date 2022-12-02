package pwes

import (
	"context"
	"testing"

	"pathwar.land/pathwar/v2/go/pkg/pwapi"

	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func TestRebuild(t *testing.T) {

	tests := []struct {
		name          string
		apiInput      *pwapi.HTTPClient
		optsInput     Opts
		expectedErr   error
		expectedScore int
	}{
		{"api null", nil, Opts{}, errcode.ErrMissingInput, 0},
	}
	var ctx context.Context
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Rebuild(ctx, test.apiInput, test.optsInput)
			assert.Equal(t, test.expectedErr, err, test.name)
		})
	}
}
