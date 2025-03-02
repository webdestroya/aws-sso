package cmd

import "github.com/webdestroya/aws-sso/internal/runners/listrunner"

func init() {
	rootCmd.AddCommand(listrunner.NewListCmd(cmdFactory))
}
