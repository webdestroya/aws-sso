package listrunner

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

func RunE(cmd *cobra.Command, args []string) error {

	return (&profileLister{cmd: cmd}).run()
}

type profileLister struct {
	cmd *cobra.Command
}

func (pl *profileLister) run() error {
	pl.cmd.Printf("Listing SSO profiles in %s:\n", appconfig.GetAwsConfigPath())
	pl.cmd.Println()

	profiles := profilepicker.Profiles()

	if len(profiles) == 0 {
		pl.cmd.PrintErr("No SSO profiles were found.")
		return nil
	}

	for _, profile := range profiles {
		if err := pl.processProfile(profile); err != nil {
			pl.cmd.Println(utils.ErrorStyle.Render("Error:", err.Error()))
			continue
		}
		pl.cmd.Println()
	}

	return nil
}

func (pl *profileLister) processProfile(profile string) error {
	pl.cmd.Printf(utils.HeaderStyle.Render("Profile: %s")+"\n", profile)

	cfg, err := awsutils.LoadSharedConfigProfile(pl.cmd.Context(), profile)
	if err != nil {
		// return utils.ErrorStyle.Render("Error:", err.Error())
		return err
	}

	ssoSession, err := awsutils.ExtractSSOInfo(cfg)
	if err != nil {
		return err
	}

	pl.cmd.Printf("Region:  %s\n", ssoSession.SSORegion)
	pl.cmd.Printf("Account: %s\n", cfg.SSOAccountID)
	pl.cmd.Printf("Role:    %s\n", cfg.SSORoleName)

	return nil
}

func oldRunE(cmd *cobra.Command, args []string) error {

	cmd.Printf("Listing SSO configurations in %s:\n", appconfig.GetAwsConfigPath())
	cmd.Println()

	entries, err := GetSSOEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		cmd.Println("No 'sso-session' entries found.")
		return nil
	}

	for _, entry := range entries {
		cmd.Printf(utils.HeaderStyle.Render("SSO: %s")+"\n", entry.ID())

		if entry.IsLegacy() {
			cmd.Println("  " + utils.WarningStyle.Render("* This profile is using legacy configuration style *"))
		}

		if len(entry.Profiles) > 0 {
			cmd.Printf("  Used By: %s\n", strings.Join(entry.Profiles, ", "))
		} else {
			cmd.Print("  Used By: Not used by any profiles\n")
		}

		if token, err := awsutils.ReadTokenInfo(entry.ID()); err == nil {

			if token.Expired() {
				cmd.Printf("  Token: %s\n", utils.WarningStyle.Render("Expired"))
			} else {
				cmd.Printf("  Token: %s, Expires: %s\n", utils.SuccessStyle.Render("Valid"), token.ExpiresAt.String())

			}

		} else {
			cmd.Println("  Token: No token file found")
		}

		cmd.Println()
	}

	return nil
}
