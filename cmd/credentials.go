package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/runners/credentialsrunner"
)

var credentialsCmd = &cobra.Command{
	Use:          "credentials PROFILE",
	Short:        "Use SSO creds as AWS process credentials",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         credentialsrunner.RunE,
}

func init() {
	rootCmd.AddCommand(credentialsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// credentialsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// credentialsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
