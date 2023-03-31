package pwapi

import (
	"github.com/stretchr/testify/assert"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"testing"
	"time"
)

func TestSeasonRules_ParseSeasonRulesString(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expectedRules SeasonRules
		expectedErr   error
	}{
		{
			name: "valid YAML",
			input: []byte(`
start_datetime: 2023-04-01T00:00:00Z
end_datetime: 2023-04-30T23:59:59Z
limit_players_per_team: 3
limit_total_teams: 5
email_domain: example.com
`),
			expectedRules: SeasonRules{
				StartDatetime:       time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDatetime:         time.Date(2023, 4, 30, 23, 59, 59, 0, time.UTC),
				LimitPlayersPerTeam: 3,
				LimitTotalTeams:     5,
				EmailDomain:         "example.com",
			},
			expectedErr: nil,
		},
		{
			name: "invalid YAML",
			input: []byte(`
				invalid
			`),
			expectedRules: SeasonRules{},
			expectedErr:   errcode.ErrParseSeasonRule,
		},
		{
			name: "start datetime after end datetime",
			input: []byte(`
start_datetime: "2023-04-30T00:00:00Z"
end_datetime: "2023-04-01T23:59:59Z"
`),
			expectedRules: SeasonRules{},
			expectedErr:   errcode.ErrSeasonRuleStartDateGreaterThanEndDate,
		},
		{
			name: "negative limit players per team",
			input: []byte(`
limit_players_per_team: -1
`),
			expectedRules: SeasonRules{},
			expectedErr:   errcode.ErrSeasonRuleInvalidLimitPlayersPerTeam,
		},
		{
			name: "negative limit total teams",
			input: []byte(`
limit_total_teams: -1
`),
			expectedRules: SeasonRules{},
			expectedErr:   errcode.ErrSeasonRuleInvalidLimitTotalTeams,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules := NewSeasonRules()
			err := rules.ParseSeasonRulesString(test.input)
			testSameErrcodes(t, "err", test.expectedErr, err)
			if err != nil {
				return
			}

			assert.Equalf(t, test.expectedRules, rules, "rules")
		})
	}
}

func TestSeasonRules_IsStarted(t *testing.T) {
	tests := []struct {
		name        string
		input       SeasonRules
		expectedRes bool
	}{
		{
			name: "started",
			input: SeasonRules{
				StartDatetime: time.Now().Add(-1 * time.Hour),
			},
			expectedRes: true,
		},
		{
			name: "not started",
			input: SeasonRules{
				StartDatetime: time.Now().Add(1 * time.Hour),
			},
			expectedRes: false,
		},
		{
			name: "no start datetime",
			input: SeasonRules{
				StartDatetime: time.Time{},
			},
			expectedRes: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.input.IsStarted()
			assert.Equalf(t, test.expectedRes, res, "res")
		})
	}
}

func TestSeasonRules_IsEnded(t *testing.T) {
	tests := []struct {
		name        string
		input       SeasonRules
		expectedRes bool
	}{
		{
			name: "ended",
			input: SeasonRules{
				EndDatetime: time.Now().Add(-1 * time.Hour),
			},
			expectedRes: true,
		},
		{
			name: "not ended",
			input: SeasonRules{
				EndDatetime: time.Now().Add(1 * time.Hour),
			},
			expectedRes: false,
		},
		{
			name: "no end datetime",
			input: SeasonRules{
				EndDatetime: time.Time{},
			},
			expectedRes: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.input.IsEnded()
			assert.Equalf(t, test.expectedRes, res, "res")
		})
	}
}

func TestSeasonRules_IsLimitTotalTeamsReached(t *testing.T) {
	tests := []struct {
		name        string
		input       SeasonRules
		teamsCount  int32
		expectedRes bool
	}{
		{
			name: "limit reached",
			input: SeasonRules{
				LimitTotalTeams: 1,
			},
			teamsCount:  1,
			expectedRes: true,
		},
		{
			name: "limit not reached",
			input: SeasonRules{
				LimitTotalTeams: 2,
			},
			teamsCount:  1,
			expectedRes: false,
		},
		{
			name: "no limit",
			input: SeasonRules{
				LimitTotalTeams: 0,
			},
			teamsCount:  1,
			expectedRes: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.input.IsLimitTotalTeamsReached(test.teamsCount)
			assert.Equalf(t, test.expectedRes, res, "res")
		})
	}
}

func TestSeasonRules_IsLimitPlayersPerTeamReached(t *testing.T) {
	tests := []struct {
		name         string
		input        SeasonRules
		playersCount int32
		expectedRes  bool
	}{
		{
			name: "limit reached",
			input: SeasonRules{
				LimitPlayersPerTeam: 1,
			},
			playersCount: 1,
			expectedRes:  true,
		},
		{
			name: "limit not reached",
			input: SeasonRules{
				LimitPlayersPerTeam: 2,
			},
			playersCount: 1,
			expectedRes:  false,
		},
		{
			name: "no limit",
			input: SeasonRules{
				LimitPlayersPerTeam: 0,
			},
			playersCount: 1,
			expectedRes:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.input.IsLimitPlayersPerTeamReached(test.playersCount)
			assert.Equalf(t, test.expectedRes, res, "res")
		})
	}
}
