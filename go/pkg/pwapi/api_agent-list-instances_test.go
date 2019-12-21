package pwapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"pathwar.land/go/internal/testutil"
	"pathwar.land/go/pkg/errcode"
)

func TestService_AgentListInstances(t *testing.T) {
	svc, cleanup := TestingService(t, ServiceOpts{Logger: testutil.Logger(t)})
	defer cleanup()
	ctx := testingSetContextToken(context.Background(), t)

	agents := testingAgents(t, svc)
	agentMap := map[string]int64{}
	for _, agent := range agents.Items {
		agentMap[agent.Name] = agent.ID
	}

	var tests = []struct {
		name        string
		input       *AgentListInstances_Input
		expectedErr error
	}{
		{"nil", nil, errcode.ErrMissingInput},
		{"empty", &AgentListInstances_Input{}, errcode.ErrMissingInput},
		{"invalid-agent", &AgentListInstances_Input{AgentID: 4242}, errcode.ErrGetAgent},
		{"localhost", &AgentListInstances_Input{AgentID: agentMap["localhost"]}, nil},
		{"localhost-2", &AgentListInstances_Input{AgentID: agentMap["localhost-2"]}, nil},
		{"inactive-agent", &AgentListInstances_Input{AgentID: agentMap["localhost-3"]}, errcode.ErrGetAgent},
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
				assert.Equal(t, test.input.AgentID, instance.AgentID)
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
