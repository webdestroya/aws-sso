package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "awssso",
	Short: "Facilitates usage of AWS sso authentication for older apps",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

var verboseLogging = false

func Execute(ver string, gitsha string) {
	rootCmd.SetVersionTemplate(`{{.Version}}`)
	rootCmd.Version = fmt.Sprintf("awssso/%s git/%s go/%s os/%s arch/%s",
		ver,
		gitsha,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	err := rootCmd.ExecuteContext(context.TODO())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.awsssogo.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verboseLogging, "verbose", false, "Verbose logging")
}
