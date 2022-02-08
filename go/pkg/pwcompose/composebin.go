package pwcompose

import (
	"os"

	"github.com/gobuffalo/packr/v2"
	"moul.io/u"
)

func tmpComposeBin() (*os.File, func(), error) {
	var pwinitBox = packr.New("binaries", "../../out")
	composeBin, err := pwinitBox.Find("docker-compose-dab")
	if err != nil {
		return nil, nil, err
	}

	file, cleanup, err := u.TempfileWithContent(composeBin)
	if err != nil {
		return nil, nil, err
	}

	if err := file.Close(); err != nil {
		return nil, nil, err
	}

	if err := os.Chmod(file.Name(), 0555); err != nil {
		return nil, nil, err
	}

	return file, cleanup, err
}
