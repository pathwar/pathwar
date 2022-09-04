package pwinit

import "github.com/gobuffalo/packr/v2"

func Binary() ([]byte, error) {
	pwinitBox := packr.New("binaries", "../../out")
	return pwinitBox.Find("pwinit-linux-amd64")
}
