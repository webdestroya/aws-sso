package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/profilelist"
	"github.com/webdestroya/aws-sso/internal/runners/credentialsrunner"
)

var credentialsCmd = &cobra.Command{
	Use:                   "credentials PROFILE",
	Short:                 "Use SSO creds as AWS process credentials",
	DisableFlagsInUseLine: true,
	ValidArgsFunction:     profilelist.ProfileCompletions,
	Args:                  cobra.MatchAll(profilelist.ValidSingleProfileArg),
	RunE:                  credentialsrunner.RunE,
	// SilenceUsage:          true,
}

func init() {
	rootCmd.AddCommand(credentialsCmd)
}
