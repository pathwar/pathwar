package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
)

func TestService_AgentRegister(t *testing.T) {
	t.Run("table", func(t *testing.T) {
		svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
		defer cleanup()
		ctx := testingSetContextToken(context.Background(), t)

		var tests = []struct {
			name        string
			input       *AgentRegister_Input
			expectedErr error
		}{
			{"nil", nil, errcode.ErrMissingInput},
			{"empty", &AgentRegister_Input{}, errcode.ErrMissingInput},
			{"new-simple", &AgentRegister_Input{Name: "just-a-test"}, nil},
			{"new-complex", &AgentRegister_Input{Name: "aaaa", Hostname: "bbbb", Arch: "cccc", OS: "dddd", Tags: []string{"eeee", "ffff"}, Version: "gggg"}, nil},
			// FIXME: check for permissions
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				ret, err := svc.AgentRegister(ctx, test.input)
				assert.Equal(t, errcode.Code(test.expectedErr), errcode.Code(err))
				if err != nil {
					return
				}

				assert.Equal(t, test.input.Name, ret.Agent.Name)
				assert.Equal(t, test.input.Hostname, ret.Agent.Hostname)
				assert.Equal(t, test.input.Arch, ret.Agent.Arch)
				assert.Equal(t, test.input.OS, ret.Agent.OS)
				assert.Equal(t, test.input.Tags, ret.Agent.TagSlice())
				assert.Equal(t, test.input.Version, ret.Agent.Version)
				assert.NotEmpty(t, ret.Agent.CreatedAt)
				assert.NotEmpty(t, ret.Agent.UpdatedAt)
				assert.NotEmpty(t, ret.Agent.LastSeenAt)
				assert.NotEmpty(t, ret.Agent.LastRegistrationAt)
			})
		}
	})
	t.Run("workflow", func(t *testing.T) {
		svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
		defer cleanup()
		ctx := testingSetContextToken(context.Background(), t)

		first, err := svc.AgentRegister(ctx, &AgentRegister_Input{Name: "test", Hostname: "lorem ipsum"})
		require.NoError(t,err)
		assert.Equal(t, first.Agent.CreatedAt, first.Agent.UpdatedAt)
		assert.Equal(t, first.Agent.LastSeenAt, first.Agent.LastRegistrationAt)
		assert.Equal(t, "lorem ipsum", first.Agent.Hostname)

		second, err := svc.AgentRegister(ctx, &AgentRegister_Input{Name: "test"})
		require.NoError(t,err)
		assert.Equal(t, first.Agent.CreatedAt, second.Agent.CreatedAt)
		assert.NotEqual(t, second.Agent.CreatedAt, second.Agent.UpdatedAt)
		assert.NotEqual(t, first.Agent.UpdatedAt, second.Agent.UpdatedAt)
		assert.Equal(t, second.Agent.LastSeenAt, second.Agent.LastRegistrationAt)
		assert.NotEqual(t, first.Agent.LastSeenAt, second.Agent.LastSeenAt)
		assert.NotEqual(t, first.Agent.LastRegistrationAt, second.Agent.LastRegistrationAt)
		assert.Empty(t, second.Agent.Hostname)
	})
}
