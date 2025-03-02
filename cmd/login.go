package cmd

import "github.com/webdestroya/aws-sso/internal/runners/loginrunner"

func init() {
	rootCmd.AddCommand(loginrunner.NewLoginCmd(cmdFactory))
}
