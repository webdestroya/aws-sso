package profilepicker

import (
	"slices"
	"strings"
	"sync"

	"github.com/webdestroya/aws-sso/internal/appconfig"
	"gopkg.in/ini.v1"
)

const (
	ssoSessionKey  = `sso_session`
	ssoStartUrlKey = `sso_start_url`
)

var profOnce = sync.OnceValue(buildProfileList)

func Profiles() []string {
	return profOnce()
}

func buildProfileList() []string {
	cfgFileIni, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		Insensitive:             true,
		AllowNestedValues:       true,
		Loose:                   true,
	}, appconfig.GetAwsConfigPath())
	if err != nil {
		return []string{}
	}

	if len(cfgFileIni.Sections()) == 0 {
		return []string{}
	}

	profiles := make([]string, 0, 10)

	for _, section := range cfgFileIni.Sections() {
		sectName := section.Name()

		if profName, has := strings.CutPrefix(sectName, "profile "); has {
			if section.HasKey(ssoSessionKey) || section.HasKey(ssoStartUrlKey) {
				profiles = append(profiles, profName)
			}
		}
	}

	slices.Sort(profiles)

	return profiles
}
