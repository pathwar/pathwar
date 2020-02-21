package pwdb

import (
	"fmt"
	"strings"
)

func newOfficialChallengeWithFlavor(name string, homepage string, composeBundle string) *ChallengeFlavor {
	return &ChallengeFlavor{
		Challenge: &Challenge{
			Name:     name,
			Author:   "Staff Team",
			Homepage: homepage,
			IsDraft:  false,
		},
		SourceURL:     homepage,
		IsLatest:      true,
		IsDraft:       false,
		Changelog:     "Initial Version",
		Version:       "v1",
		ComposeBundle: composeBundle,
		Driver:        ChallengeFlavor_DockerCompose,
	}
}

func (cf *ChallengeFlavor) addSeasonChallengeByID(seasonID int64) {
	if cf.SeasonChallenges == nil {
		cf.SeasonChallenges = []*SeasonChallenge{}
	}
	cf.SeasonChallenges = append(cf.SeasonChallenges, &SeasonChallenge{
		SeasonID: seasonID,
	})
}

func (a *Agent) TagSlice() []string {
	if a.Tags == "" {
		return nil
	}
	return strings.Split(a.Tags, ", ")
}

func (cf ChallengeFlavor) NameAndVersion() string {
	return fmt.Sprintf("%s@%s", cf.Challenge.Name, cf.Version)
}
