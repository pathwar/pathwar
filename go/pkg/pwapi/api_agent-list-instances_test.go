package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/v2/go/internal/testutil"
	"pathwar.land/v2/go/pkg/errcode"
)

func TestService_AgentListInstances(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	var tests = []struct {
		name        string
		input       *AgentListInstances_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &AgentListInstances_Input{}, errcode.ErrMissingInput},
		{"invalid-agent", &AgentListInstances_Input{AgentName: "unknown"}, errcode.ErrGetAgent},
		{"localhost", &AgentListInstances_Input{AgentName: "localhost"}, nil},
		{"localhost-2", &AgentListInstances_Input{AgentName: "localhost-2"}, nil},
		{"inactive-agent", &AgentListInstances_Input{AgentName: "localhost-3"}, errcode.ErrGetAgent},
		// FIXME: check for permissions
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := svc.AgentListInstances(ctx, test.input)
			testSameErrcodes(t, "", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Len(t, ret.Instances, 2) // FIXME: should be 1 if we keep only active ones
			for _, instance := range ret.Instances {
				assert.NotNil(t, instance.Agent)
				assert.Equal(t, test.input.AgentName, instance.Agent.Name)
				assert.NotNil(t, instance.Flavor)
				assert.NotNil(t, instance.Flavor.Challenge)
				for _, seasonChallenge := range instance.Flavor.SeasonChallenges {
					assert.Equal(t, instance.Flavor.ID, seasonChallenge.FlavorID)
					// FIXME: verify seasonChallenge.ChallengeSubscriptions...
				}
			}
			//fmt.Println(godev.PrettyJSON(ret))
		})
	}
}
