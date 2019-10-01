package pwlevel

import "github.com/gobuffalo/packr/v2"

func Binary() ([]byte, error) {
	var pwlevelBox = packr.New("pwctl-binaries", "../../out")
	return pwlevelBox.Find("pwlevel-linux-amd64")
}
