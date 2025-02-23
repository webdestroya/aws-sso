package listrunner

/*
import (
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"gopkg.in/ini.v1"
)

// Deprecated
// use profilelist.ProfileCompletions() instead
func ProfileCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	// TODO: use profilelist.Profiles

	cfgFileIni, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		Insensitive:             true,
		AllowNestedValues:       true,
		Loose:                   true,
	}, appconfig.GetAwsConfigPath())
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	if len(cfgFileIni.Sections()) == 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	profiles := make([]string, 0, 10)

	toComplete = strings.ToLower(toComplete)

	// fmt.Printf("COMPLETE: [%s]\n\n", toComplete)

	for _, section := range cfgFileIni.Sections() {
		sectName := section.Name()

		if profName, has := strings.CutPrefix(sectName, "profile "); has {
			if section.HasKey(ssoSessionKey) || section.HasKey(ssoStartUrlKey) {

				if toComplete != "" && !strings.HasPrefix(profName, toComplete) {
					continue
				}

				profiles = append(profiles, profName)
			}
		}
	}

	slices.Sort(profiles)

	return profiles, cobra.ShellCompDirectiveNoFileComp
}
*/
