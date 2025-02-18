package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/runners/envrunner"
)

var envCmd = &cobra.Command{
	Use:                   "env PROFILE -- command [command-args...]",
	Aliases:               []string{"run"},
	Short:                 "Run a command with AWS access keys set in the environment",
	SilenceUsage:          false,
	DisableFlagsInUseLine: true,
	FParseErrWhitelist:    cobra.FParseErrWhitelist{UnknownFlags: true},
	Args:                  cobra.MinimumNArgs(2),
	RunE:                  envrunner.RunE,
}

func init() {
	rootCmd.AddCommand(envCmd)
}
