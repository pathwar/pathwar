package pwapi

import (
	"strings"
	"time"

	"pathwar.land/pathwar/v2/go/pkg/errcode"

	"gopkg.in/yaml.v3"
)

type SeasonRules struct {
	StartDatetime       time.Time `yaml:"start_datetime"`
	EndDatetime         time.Time `yaml:"end_datetime"`
	LimitPlayersPerTeam int32     `yaml:"limit_players_per_team"`
	LimitTotalTeams     int32     `yaml:"limit_total_teams"`
	EmailDomain         string    `yaml:"email_domain"`
}

func (s *SeasonRules) ParseSeasonRulesString(seasonsRulesYAML []byte) error {
	err := yaml.Unmarshal(seasonsRulesYAML, s)
	if err != nil {
		return errcode.ErrParseSeasonRule
	}
	if s.StartDatetime.Unix() > 0 && s.EndDatetime.Unix() > 0 && s.StartDatetime.Unix() > s.EndDatetime.Unix() {
		return errcode.ErrSeasonRuleStartDateGreaterThanEndDate
	}
	if s.LimitPlayersPerTeam < 0 {
		return errcode.ErrSeasonRuleInvalidLimitPlayersPerTeam
	}
	if s.LimitTotalTeams < 0 {
		return errcode.ErrSeasonRuleInvalidLimitTotalTeams
	}
	return nil
}

func (s *SeasonRules) IsStarted() bool {
	return s.StartDatetime.Unix() > 0 && s.StartDatetime.Unix() <= time.Now().Unix()
}

func (s *SeasonRules) IsEnded() bool {
	return s.EndDatetime.Unix() > 0 && s.EndDatetime.Unix() <= time.Now().Unix()
}

func (s *SeasonRules) LimitTotalTeamsReached(totalTeams int32) bool {
	return s.LimitTotalTeams > 0 && totalTeams >= s.LimitTotalTeams
}

func (s *SeasonRules) LimitPlayersPerTeamReached(totalPlayers int32) bool {
	return s.LimitPlayersPerTeam > 0 && totalPlayers >= s.LimitPlayersPerTeam
}

func (s *SeasonRules) IsEmailDomainAllowed(email string) bool {
	domain := email[strings.LastIndex(email, "@")+1:]
	return s.EmailDomain == "" || domain == s.EmailDomain
}

func NewSeasonRules() SeasonRules {
	return SeasonRules{
		StartDatetime:       time.Now(),
		EndDatetime:         time.Time{},
		LimitPlayersPerTeam: 0,
		LimitTotalTeams:     0,
		EmailDomain:         "",
	}
}
