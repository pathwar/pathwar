package pwdb

func newOfficialChallengeWithFlavor(name string, homepage string) *ChallengeFlavor {
	return &ChallengeFlavor{
		Challenge: &Challenge{
			Name:     name,
			Author:   "Staff Team",
			Homepage: homepage,
			IsDraft:  false,
		},
		SourceURL: homepage,
		IsLatest:  true,
		IsDraft:   false,
		Changelog: "Initial Version",
		Version:   "v1",
		Driver:    ChallengeFlavor_DockerCompose,
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
