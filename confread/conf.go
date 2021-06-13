package confread

import (
	"errors"
	"io/ioutil"

	"github.com/bagaking/gotools/file/fpth"
)

var ErrNotFound = errors.New("config file not found")

func Read(fileName string, out interface{}, positions ...string) error {
	allPath := make([]string, 0)
	for _, pos := range positions {
		dir, err := fpth.Adapt(pos, fpth.OEnableHomeDir(), fpth.ORelativePWDPath())
		if err != nil {
			return err
		}

		allPath = append(allPath, fpth.Join(dir, fileName))
	}

	var data []byte = nil
	for _, confPth := range allPath {
		result, err := ioutil.ReadFile(confPth)
		if err == nil {
			data = result
			break
		}
	}

	if data == nil {
		return ErrNotFound
	}

	return TryAllUnmarshalers(data, out)
}
