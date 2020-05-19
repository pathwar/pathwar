package pwdb

import (
	"fmt"
	"strings"

	"github.com/martinlindhe/base36"
	"golang.org/x/crypto/sha3"
	"pathwar.land/v2/go/pkg/errcode"
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

func ChallengeInstancePrefixHash(instanceID int64, userID int64, salt string) (string, error) {
	stringToHash := fmt.Sprintf("%s%d%s", fmt.Sprintf("%d", instanceID), userID, salt)
	hashBytes := make([]byte, 8)
	hasher := sha3.NewShake256()
	_, err := hasher.Write([]byte(stringToHash))
	if err != nil {
		return "", errcode.ErrWriteBytesToHashBuilder.Wrap(err)
	}
	_, err = hasher.Read(hashBytes)
	if err != nil {
		return "", errcode.ErrReadBytesFromHashBuilder.Wrap(err)
	}
	userHash := strings.ToLower(base36.EncodeBytes(hashBytes))[:8] // we voluntarily expect short hashes here
	return userHash, nil
}
