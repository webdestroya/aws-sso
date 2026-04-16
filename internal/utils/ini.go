package utils

import (
	"errors"
	"os"

	"github.com/webdestroya/aws-sso/internal/utils/fsutils"
	"gopkg.in/ini.v1"
)

func LoadIniFile(opts ini.LoadOptions, filename string) (*ini.File, error) {

	data, err := fsutils.ReadFile(filename)
	if err != nil {
		if opts.Loose && errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
	}

	return ini.LoadSources(opts, data)
}
