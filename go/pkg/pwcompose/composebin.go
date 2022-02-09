package pwcompose

import (
	"os"

	"github.com/gobuffalo/packr/v2"
	"moul.io/u"
)

func ComposeBinBytes() ([]byte, error) {
	var pwinitBox = packr.New("binaries", "../../out")
	return pwinitBox.Find("docker-compose-dab")
}

func tmpComposeBin() (*os.File, func(), error) {
	composeBin, err := ComposeBinBytes()
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
