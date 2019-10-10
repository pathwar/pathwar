package pwchallenge

import "github.com/gobuffalo/packr/v2"

func Binary() ([]byte, error) {
	var pwchallengeBox = packr.New("pwctl-binaries", "../../out")
	return pwchallengeBox.Find("pwchallenge-linux-amd64")
}
