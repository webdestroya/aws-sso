package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:                   "login PROFILE",
	Short:                 "Login to the SSO session for the specified profile",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE:                  runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

// https://github.com/synfinatic/aws-sso-cli/blob/main/internal/sso/awssso.go
// https://github.com/aws/aws-cli/blob/v2/awscli/customizations/sso/utils.py

func runLogin(cmd *cobra.Command, args []string) error {

	awsCliPath, err := exec.LookPath("aws")
	if err != nil {
		return err
	}

	profileName := args[0]
	cmdArgs := []string{
		"aws", "sso", "login", "--profile", profileName,
	}

	// cmd.Printf("Running: %s %v\n\n", awsCliPath, cmdArgs)

	envVars := os.Environ()
	return syscall.Exec(awsCliPath, cmdArgs, envVars)
}

func doSsoLogin() error {
	return nil
}
