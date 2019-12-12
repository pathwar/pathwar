package pwapi

import (
	"context"
	"testing"

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
				testSameErrcodes(t, "", test.expectedErr, err)
				if err != nil {
					return
				}

				testSameStrings(t, "", test.input.Name, ret.Agent.Name)
				testSameStrings(t, "", test.input.Hostname, ret.Agent.Hostname)
				testSameStrings(t, "", test.input.Arch, ret.Agent.Arch)
				testSameStrings(t, "", test.input.OS, ret.Agent.OS)
				testSameDeep(t, "", test.input.Tags, ret.Agent.TagSlice())
				testSameStrings(t, "", test.input.Version, ret.Agent.Version)
				testIsTrue(t, "", ret.Agent.CreatedAt != nil && !ret.Agent.CreatedAt.IsZero())
				testIsTrue(t, "", ret.Agent.UpdatedAt != nil && !ret.Agent.UpdatedAt.IsZero())
				testIsTrue(t, "", ret.Agent.LastSeenAt != nil && !ret.Agent.LastSeenAt.IsZero())
				testIsTrue(t, "", ret.Agent.LastRegistrationAt != nil && !ret.Agent.LastRegistrationAt.IsZero())
			})
		}
	})
	t.Run("workflow", func(t *testing.T) {
		svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
		defer cleanup()
		ctx := testingSetContextToken(context.Background(), t)

		first, err := svc.AgentRegister(ctx, &AgentRegister_Input{Name: "test", Hostname: "lorem ipsum"})
		checkErr(t, "", err)
		testIsTrue(t, "", first.Agent.CreatedAt.Equal(*first.Agent.UpdatedAt))
		testIsTrue(t, "", first.Agent.LastSeenAt.Equal(*first.Agent.LastRegistrationAt))
		testSameStrings(t, "", "lorem ipsum", first.Agent.Hostname)

		second, err := svc.AgentRegister(ctx, &AgentRegister_Input{Name: "test"})
		checkErr(t, "", err)
		testIsTrue(t, "", first.Agent.CreatedAt.Equal(*second.Agent.CreatedAt))
		testIsTrue(t, "", !second.Agent.CreatedAt.Equal(*second.Agent.UpdatedAt))
		testIsTrue(t, "", !first.Agent.UpdatedAt.Equal(*second.Agent.UpdatedAt))
		testIsTrue(t, "", second.Agent.LastSeenAt.Equal(*second.Agent.LastRegistrationAt))
		testIsTrue(t, "", !first.Agent.LastSeenAt.Equal(*second.Agent.LastSeenAt))
		testIsTrue(t, "", !first.Agent.LastRegistrationAt.Equal(*second.Agent.LastRegistrationAt))
		testSameStrings(t, "", "", second.Agent.Hostname)
	})
}
