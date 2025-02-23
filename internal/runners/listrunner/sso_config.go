package listrunner

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/webdestroya/aws-sso/internal/appconfig"
	"gopkg.in/ini.v1"
)

// var (
// 	ssoSessionRegexp = regexp.MustCompile(`^\[sso-session ([-_a-zA-Z0-9]+)\]`)
// )

const (
	ssoSessionKey  = `sso_session`
	ssoStartUrlKey = `sso_start_url`
)

type SSOEntry struct {
	Name     string
	StartURL string
	Profiles []string
}

func (s SSOEntry) IsLegacy() bool {
	return s.Name == ""
}

func (s SSOEntry) String() string {
	if s.IsLegacy() {
		return fmt.Sprintf("Legacy: %s", s.StartURL)
	}
	return s.Name
}

func (s SSOEntry) ID() string {
	if s.IsLegacy() {
		return s.StartURL
	}
	return s.Name
}

type ssoEntryBuilder struct {
	ssodb map[string]*SSOEntry
}

func GetSSOEntries() ([]SSOEntry, error) {
	b := &ssoEntryBuilder{
		ssodb: make(map[string]*SSOEntry),
	}
	return b.buildMap()
}

func (b *ssoEntryBuilder) ensureEntry(key string) {
	if _, ok := b.ssodb[key]; ok {
		return
	}

	b.ssodb[key] = &SSOEntry{
		Profiles: make([]string, 0, 10),
	}
}

func (b *ssoEntryBuilder) setName(key, value string) {
	b.ensureEntry(key)
	if b.ssodb[key].Name != "" {
		return
	}
	b.ssodb[key].Name = value
}

func (b *ssoEntryBuilder) setUrl(key, value string) {
	b.ensureEntry(key)
	if b.ssodb[key].StartURL != "" {
		return
	}
	b.ssodb[key].StartURL = value
}

func (b *ssoEntryBuilder) addProfile(key, profile string) {
	b.ensureEntry(key)
	newProfiles := append(b.ssodb[key].Profiles, profile)
	slices.Sort(newProfiles)
	b.ssodb[key].Profiles = slices.Compact(newProfiles)
}

func (b *ssoEntryBuilder) buildMap() ([]SSOEntry, error) {
	cfgFileIni, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		Insensitive:             true,
		AllowNestedValues:       true,
		Loose:                   true,
	}, appconfig.GetAwsConfigPath())
	if err != nil {
		return nil, err
	}

	if len(cfgFileIni.Sections()) == 0 {
		// cmd.Println(utils.ErrorStyle.Render("ERROR"), "Failed to read/parse config file", err.Error())
		return nil, nil
	}

	for _, section := range cfgFileIni.Sections() {
		sectName := section.Name()

		if ssoName, has := strings.CutPrefix(sectName, "sso-session "); has {

			b.ensureEntry(ssoName)
			b.setName(ssoName, ssoName)
			if kv, err := section.GetKey(ssoStartUrlKey); err == nil {
				b.setUrl(ssoName, kv.MustString(""))
			}

		} else if profName, has := strings.CutPrefix(sectName, "profile "); has {
			if skey, serr := section.GetKey(ssoSessionKey); serr == nil {
				ssoName := skey.String()
				b.addProfile(ssoName, profName)
				b.setName(ssoName, ssoName)

			} else if skey, serr := section.GetKey(ssoStartUrlKey); serr == nil {
				ssoUrl := skey.String()
				b.addProfile(ssoUrl, profName)
				b.setUrl(ssoUrl, ssoUrl)
			}
		}
	}

	if len(b.ssodb) == 0 {
		return nil, nil
	}

	keys := slices.Collect(maps.Keys(b.ssodb))
	slices.Sort(keys)

	ssoList := make([]SSOEntry, 0, len(b.ssodb))
	for _, key := range keys {
		ssoList = append(ssoList, *b.ssodb[key])
	}

	return ssoList, nil
}
