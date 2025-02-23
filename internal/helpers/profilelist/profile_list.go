package profilelist

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cobra"
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

func ValidSingleProfileArg(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no profile argument was provided")
	}
	if len(args) > 1 {
		return fmt.Errorf("only one profile can be provided")
	}

	profiles := Profiles()
	profile := args[0]
	if !slices.Contains(profiles, profile) {
		return fmt.Errorf("invalid argument: %s is not an SSO profile", profile)
	}

	return nil
}

func ValidProfileArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}

	profiles := Profiles()

	for _, profile := range args {
		if !slices.Contains(profiles, profile) {
			return fmt.Errorf("invalid argument: %s is not an SSO profile", profile)
		}
	}
	return nil
}

func ProfileCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	possibleProfiles := Profiles()

	if len(possibleProfiles) == 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	completions := make([]string, 0, len(possibleProfiles))
	toComplete = strings.ToLower(toComplete)
	for _, profile := range possibleProfiles {

		if toComplete != "" && !strings.HasPrefix(profile, toComplete) {
			continue
		}

		completions = append(completions, profile)
	}

	slices.Sort(completions)

	return completions, cobra.ShellCompDirectiveNoFileComp
}
