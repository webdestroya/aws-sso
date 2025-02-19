package doctorrunner

import (
	"errors"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/utils"
	"gopkg.in/ini.v1"
)

var (
	ssoSessionRegexp = regexp.MustCompile(`^\[sso-session ([-_a-zA-Z0-9]+)\]`)
)

const (
	ssoSessionKey  = `sso_session`
	ssoStartUrlKey = `sso_start_url`
)

func RunE(cmd *cobra.Command, args []string) error {
	if err := checkAwsConfig(cmd); err != nil {
		return err
	}

	return nil
}

func checkAwsConfig(cmd *cobra.Command) error {

	awsCfgFile := config.DefaultSharedConfigFilename()

	cmd.Printf("Checking %s file...", awsCfgFile)
	if _, err := os.Stat(awsCfgFile); err == nil {
		cmd.Printf("EXISTS (%s)\n", awsCfgFile)
	} else if errors.Is(err, os.ErrNotExist) {
		cmd.Println(utils.ErrorStyle.Render("MISSING"))
		cmd.Println("Skipping configuration checks!")
		return nil
	} else {
		cmd.Println(utils.ErrorStyle.Render("ERROR"), err.Error())
		cmd.Println("Skipping configuration checks!")
		return nil
	}

	// cfgBytes, err := os.ReadFile(awsCfgFile)
	// if err != nil {
	// 	cmd.Println(utils.ErrorStyle.Render("ERROR"), "Failed to read config file", err.Error())
	// }

	// cmd.Print("Checking for sso configurations...")
	// results := ssoSessionRegexp.FindAllStringSubmatch(string(cfgBytes), -1)
	// ssoNames := make([]string, 0, len(results))
	// for _, v := range results {
	// 	ssoNames = append(ssoNames, v[1])
	// }

	// slices.Sort(ssoNames)

	// if len(ssoNames) > 0 {
	// 	cmd.Printf("FOUND (%d)\n", len(ssoNames))
	// 	for _, v := range ssoNames {

	// 		cmd.Printf(" * %s\n", v)
	// 	}
	// } else {
	// 	cmd.Println(utils.ErrorStyle.Render("NONE"), "No 'sso-session' entries found. You need to configure SSO!")
	// 	return nil
	// }

	configFiles := make([]string, 0, 10)
	configFiles = append(configFiles, config.DefaultSharedConfigFiles...)
	// configFiles = append(configFiles, config.DefaultSharedCredentialsFiles...)

	configFile, otherCfgFiles := configFiles[0], configFiles[1:]

	cfgFileIni, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		Insensitive:             true,
		AllowNestedValues:       true,
		Loose:                   true,
	}, configFile, utils.ToAnySlice(otherCfgFiles)...)
	if err != nil || len(cfgFileIni.Sections()) == 0 {
		cmd.Println(utils.ErrorStyle.Render("ERROR"), "Failed to read/parse config file", err.Error())
	}

	// map of sso sessions and a list of profiles using that
	usageMap := make(map[string][]string, 0)
	legacyUsageMap := make(map[string][]string, 0)

	cmd.Print("Checking for sso configurations...")
	ssoNames := make([]string, 0, 10)
	for _, section := range cfgFileIni.Sections() {
		sectName := section.Name()

		if ssoName, has := strings.CutPrefix(sectName, "sso-session "); has {
			ssoNames = append(ssoNames, ssoName)

		} else if profName, has := strings.CutPrefix(sectName, "profile "); has {
			if section.HasKey(ssoSessionKey) {
				if skey, serr := section.GetKey(ssoSessionKey); serr == nil {
					ssoName := skey.String()
					usageMap[ssoName] = append(usageMap[ssoName], profName)
				}
			} else if section.HasKey(ssoStartUrlKey) {
				// Legacy SSO
				if skey, serr := section.GetKey(ssoStartUrlKey); serr == nil {
					ssoName := skey.String()
					legacyUsageMap[ssoName] = append(legacyUsageMap[ssoName], profName)
				}
			}
		}
	}

	foundSSO := false

	slices.Sort(ssoNames)
	if len(ssoNames) > 0 {
		cmd.Printf("FOUND (%d)\n", len(ssoNames))
		for _, v := range ssoNames {
			cmd.Printf(" * %s ", v)
			if profiles, ok := usageMap[v]; ok {
				slices.Sort(profiles)
				cmd.Printf("(Used By: %s)", strings.Join(profiles, ", "))
			} else {
				cmd.Print("(Not used by any profiles)")
			}
			cmd.Println()
		}
		foundSSO = true
	}

	// LOOK FOR LEGACY PROFILES
	legacyUrls := slices.Collect(maps.Keys(legacyUsageMap))
	if len(legacyUrls) > 0 {
		cmd.Println()
		cmd.Printf("Also found %d legacy SSO configurations:\n", len(legacyUrls))
		for _, v := range legacyUrls {
			cmd.Printf(" * %s ", v)
			if profiles, ok := legacyUsageMap[v]; ok {
				slices.Sort(profiles)
				cmd.Printf("(Used By: %s)", strings.Join(profiles, ", "))
			} else {
				cmd.Print("(Not used by any profiles)")
			}
			cmd.Println()
		}
		foundSSO = true
	}

	if !foundSSO {
		cmd.Println(utils.ErrorStyle.Render("NONE"), "No 'sso-session' entries found. You need to configure SSO!")
		return nil
	}

	return nil
}
