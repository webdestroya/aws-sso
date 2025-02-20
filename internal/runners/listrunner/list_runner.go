package listrunner

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

func RunE(cmd *cobra.Command, args []string) error {

	entries, err := GetSSOEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		cmd.Println("No 'sso-session' entries found.")
		return nil
	}

	for _, entry := range entries {
		cmd.Printf("SSO: %s\n", entry.ID())

		if len(entry.Profiles) > 0 {
			cmd.Printf("  Used By: %s\n", strings.Join(entry.Profiles, ", "))
		} else {
			cmd.Print("  Used By: Not used by any profiles\n")
		}

		if token, err := awsutils.ReadTokenInfo(entry.StartURL); err == nil {

			if token.Expired() {
				cmd.Printf("  Token: %s\n", utils.WarningStyle.Render("Expired"))
			} else {
				cmd.Printf("  Token: %s, expires: %s\n", utils.SuccessStyle.Render("Valid"), token.ExpiresAt.String())

			}

		} else {
			cmd.Println("  Token: No token file found")
		}

		cmd.Println()
	}

	return nil
}
