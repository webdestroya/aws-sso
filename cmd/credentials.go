package cmd

import "github.com/webdestroya/aws-sso/internal/runners/credentialsrunner"

func init() {
	rootCmd.AddCommand(credentialsrunner.NewCredentialsCmd(cmdFactory))
}
