//go:build !nodev

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/runners/devrunner"
)

var devCmd = &cobra.Command{
	Use:                   "dev",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	RunE:                  devrunner.RunE,
}

func init() {
	rootCmd.AddCommand(devCmd)
}
