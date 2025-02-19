package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/runners/syncrunner"
)

var syncCmd = &cobra.Command{
	Use:          "sync PROFILE [PROFILE...]",
	Short:        "Sync AWS credentials. (This will overwrite the profile credentials!)",
	Example:      "awssso sync mycompany-production",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE:         syncrunner.RunE,
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().Bool("force", false, "Force overwrite profile credentials even if they do not appear to be for an SSO profile")
}
