package pwapi

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type SeasonsRules struct {
	StartDatetime       int64  `yaml:"start_datetime"`
	EndDatetime         int64  `yaml:"end_datetime"`
	LimitPlayersPerTeam int32  `yaml:"limit_players_per_team"`
	LimitTotalTeams     int32  `yaml:"limit_total_teams"`
	EmailDomain         string `yaml:"email_domain"`
}

func readSeasonsRulesFile(seasonsRulesFilePath string) (SeasonsRules, error) {
	if seasonsRulesFilePath == "" {
		return NewSeasonsRules(), nil
	}
	seasonsRulesYAML, err := os.ReadFile(seasonsRulesFilePath)
	if err != nil {
		return SeasonsRules{}, err
	}
	return parseSeasonsRules(string(seasonsRulesYAML))
}

func parseSeasonsRules(seasonsRulesYAML string) (SeasonsRules, error) {
	seasonsRules := NewSeasonsRules()
	if err := yaml.Unmarshal([]byte(seasonsRulesYAML), &seasonsRules); err != nil {
		return seasonsRules, err
	}
	return seasonsRules, nil
}

func (s SeasonsRules) IsOpen() bool {
	return s.IsStarted() && !s.IsEnded()
}

func (s SeasonsRules) IsStarted() bool {
	return s.StartDatetime <= time.Now().Unix()
}

func (s SeasonsRules) IsEnded() bool {
	return s.EndDatetime > 0 && s.EndDatetime <= time.Now().Unix()
}

func NewSeasonsRules() SeasonsRules {
	return SeasonsRules{
		StartDatetime:       time.Now().Unix(),
		EndDatetime:         0,
		LimitPlayersPerTeam: 0,
		LimitTotalTeams:     0,
		EmailDomain:         "",
	}
}
