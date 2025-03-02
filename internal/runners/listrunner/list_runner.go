package listrunner

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/helpers/getcreds"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

type profileLister struct {
	cmd  *cobra.Command
	opts *listOptions
}

func (pl *profileLister) run() error {

	if pl.opts.ShowSessions {
		pl.listSessions()

		pl.cmd.Println(strings.Repeat("-", 80))
		pl.cmd.Println()
	}

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

	out := pl.cmd.OutOrStdout()

	fmt.Fprintln(out, utils.HeaderStyle.Render("Profile:", profile))

	cfg, err := awsutils.LoadSharedConfigProfile(pl.cmd.Context(), profile)
	if err != nil {
		// return utils.ErrorStyle.Render("Error:", err.Error())
		return err
	}

	ssoSession, err := awsutils.ExtractSSOInfo(cfg)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "  Region:    %s\n", ssoSession.SSORegion)
	fmt.Fprintf(out, "  Account:   %s\n", cfg.SSOAccountID)
	fmt.Fprintf(out, "  Role:      %s\n", cfg.SSORoleName)
	if ssoSession.Name != "" {
		fmt.Fprintf(out, "  Session:   %s\n", ssoSession.Name)
	}

	cachedCreds, _ := getcreds.ReadCliCache(cfg, ssoSession)
	if cachedCreds != nil {
		renderCacheCreds(out, cachedCreds)
	}

	return nil
}

func renderCacheCreds(out io.Writer, creds *aws.Credentials) {
	defer fmt.Fprintln(out)
	fmt.Fprint(out, "  Cached:    ")

	if creds.Expired() {
		fmt.Fprintf(out, "Expired")
	} else {
		fmt.Fprint(out, "Valid")
		if creds.CanExpire {
			fmt.Fprintf(out, " (expires %s)", creds.Expires.Local().Format(time.DateTime))
		}
	}

}

func (pl *profileLister) listSessions() error {

	out := pl.cmd.OutOrStdout()

	fmt.Fprintf(out, "Listing SSO Sessions in %s:\n", appconfig.GetAwsConfigPath())
	fmt.Fprintln(out)

	entries, err := GetSSOEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Fprintln(out, "No 'sso-session' entries found.")
		return nil
	}

	for _, entry := range entries {
		fmt.Fprintln(out, utils.HeaderStyle.Render("Session:", entry.ID()))

		if entry.IsLegacy() {
			fmt.Fprintln(out, "  "+utils.WarningStyle.Render("* This profile is using legacy configuration style *"))
		}

		fmt.Fprintln(out, "  Start URL:", entry.StartURL)

		if len(entry.Profiles) > 0 {
			fmt.Fprintln(out, "  Used By:", strings.Join(entry.Profiles, ", "))
		} else {
			fmt.Fprintln(out, "  Used By: Not used by any profiles")
		}

		if token, err := awsutils.ReadTokenInfo(entry.ID()); err == nil {

			if token.Expired() {
				fmt.Fprintln(out, "  Token:", utils.WarningStyle.Render("Expired"))
			} else {
				fmt.Fprintf(out, "  Token: %s, Expires: %s\n", utils.SuccessStyle.Render("Valid"), token.ExpiresAt.AsTime().Format(time.DateTime))

			}

		} else {
			fmt.Fprintln(out, "  Token: No token file found")
		}

		fmt.Fprintln(out)
	}

	return nil
}
