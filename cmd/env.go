package cmd

import "github.com/webdestroya/aws-sso/internal/runners/envrunner"

func init() {
	rootCmd.AddCommand(envrunner.NewEnvCmd(cmdFactory))
}
