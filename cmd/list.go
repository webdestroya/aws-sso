package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/runners/listrunner"
)

var listCmd = &cobra.Command{
	Use:                   "list",
	Short:                 "Lists available SSO sessions and any info about them",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MatchAll(cobra.NoArgs),
	ValidArgsFunction:     cobra.NoFileCompletions,
	RunE:                  listrunner.RunE,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
