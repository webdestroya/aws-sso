package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/runners/loginrunner"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:                   "login profile [profile]...",
	Short:                 "Login to the SSO session for the specified profile(s)",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE:                  loginrunner.RunE,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

// https://github.com/synfinatic/aws-sso-cli/blob/main/internal/sso/awssso.go
// https://github.com/aws/aws-cli/blob/v2/awscli/customizations/sso/utils.py
