// This will perform the login function to a specific StartURL SSO endpoint
// it will then update the cached token for the StartURL and/or Name
// it then exits, it does not actually get role credentials
package loginrunner

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/loginflow"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

func RunE(opts *LoginOptions, cmd *cobra.Command, args []string) error {
	profiles, err := profilepicker.GetProfilesFromArgsOrPrompt(cmd, args)
	if err != nil {
		return err
	}

	return runProfiles(opts, cmd, profiles)
}

func runProfiles(opts *LoginOptions, cmd *cobra.Command, profiles []string) error {
	ctx := cmd.Context()

	cfgmap := make(map[string]*config.SSOSession)

	for _, profile := range profiles {
		sharedCfg, err := awsutils.LoadSharedConfigProfile(ctx, profile)
		if err != nil {
			return err
		}

		ssoSession, err := awsutils.ExtractSSOInfo(sharedCfg)
		if err != nil {
			return err
		}

		cachePath, err := awsutils.GetSSOCachePath(ssoSession)
		if err != nil {
			return err
		}

		cfgmap[cachePath] = ssoSession
	}

	lFlowOpts := []loginflow.LoginFlowOption{
		loginflow.WithBrowser(!opts.NoBrowser),
	}

	for _, session := range cfgmap {
		if _, err := loginflow.DoLoginFlow(ctx, cmd.ErrOrStderr(), session, lFlowOpts...); err != nil {
			return err
		}
	}

	return nil
}
