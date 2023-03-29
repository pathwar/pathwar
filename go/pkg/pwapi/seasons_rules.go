package pwapi

import "time"

type SeasonsRules struct {
	StartDatetime       int64  `yaml:"start_datetime"`
	EndDatetime         int64  `yaml:"end_datetime"`
	LimitPlayersPerTeam int32  `yaml:"limit_players_per_team"`
	LimitTotalTeams     int32  `yaml:"limit_total_teams"`
	EmailDomain         string `yaml:"email_domain"`
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
