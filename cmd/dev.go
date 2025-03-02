//go:build !nodev

package cmd

import "github.com/webdestroya/aws-sso/internal/runners/devrunner"

func init() {
	rootCmd.AddCommand(devrunner.NewDevCmd(cmdFactory))
}
