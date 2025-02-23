package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
	"github.com/webdestroya/aws-sso/internal/runners/loginrunner"
	"github.com/webdestroya/aws-sso/internal/runners/syncrunner"
)

var rootCmd = &cobra.Command{
	Use:   "awssso",
	Short: "Facilitates usage of AWS sso authentication for older apps",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func Execute(ver string, gitsha string) {
	rootCmd.SetVersionTemplate(`{{.Version}}`)
	rootCmd.Version = fmt.Sprintf("awssso/%s go/%s os/%s arch/%s",
		ver,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetIn(os.Stdin)

	err := rootCmd.ExecuteContext(context.TODO())

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(rootCmd.OutOrStdout(), "Exiting..")
			return
		}

		// fmt.Fprintf(rootCmd.OutOrStderr(), "ERR=%T/msg=%v\n", err, err.Error())
		os.Exit(1)
	}
}

func init() {
	f := factory.Default()

	rootCmd.AddCommand(loginrunner.NewLoginCmd(f))
	rootCmd.AddCommand(syncrunner.NewCmdSync(f))
}
