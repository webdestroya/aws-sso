package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/runners/doctorrunner"
)

var doctorCmd = &cobra.Command{
	Use:                   "doctor",
	Short:                 "Checks for possible issues using this command",
	Args:                  cobra.NoArgs,
	DisableFlagsInUseLine: true,
	RunE:                  doctorrunner.RunE,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
