package cmd

import "github.com/webdestroya/aws-sso/internal/runners/syncrunner"

func init() {
	rootCmd.AddCommand(syncrunner.NewSyncCmd(cmdFactory))
}
