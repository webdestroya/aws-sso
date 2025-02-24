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
	"github.com/webdestroya/aws-sso/internal/runners/credentialsrunner"
	"github.com/webdestroya/aws-sso/internal/runners/loginrunner"
	"github.com/webdestroya/aws-sso/internal/runners/syncrunner"
	"github.com/webdestroya/aws-sso/internal/utils"
)

var rootCmd = &cobra.Command{
	Use:           "awssso",
	Short:         "Facilitates usage of AWS SSO authentication for older apps",
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	// DisableFlagsInUseLine: true,
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	cmd.Println("YARRR ROOT")
	// 	return pflag.ErrHelp
	// },
}

func Execute(ver string, gitsha string) int {
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

	cmd, err := rootCmd.ExecuteContextC(context.TODO())

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(rootCmd.OutOrStdout(), "Exiting..")
			return 0
		}

		// fmt.Fprintf(rootCmd.OutOrStderr(), "cmd=%v ERR=%T/msg=%v\n", cmd.Name(), err, err.Error())

		cmd.PrintErrln(utils.ErrorStyle.Render(cmd.ErrPrefix(), err.Error()))
		cmd.PrintErrf("Run '%v --help' for usage.\n", cmd.CommandPath())

		// os.Exit(1)
		return 1
	}

	return 0
}

func init() {
	f := factory.Default()

	rootCmd.AddCommand(loginrunner.NewLoginCmd(f))
	rootCmd.AddCommand(syncrunner.NewCmdSync(f))
	rootCmd.AddCommand(credentialsrunner.NewCredentialsCmd(f))
}
