package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

var cmdFactory = factory.Default()

var rootCmd = &cobra.Command{
	Use:           "awssso",
	Short:         "Facilitates usage of AWS SSO authentication for older apps",
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// if we got this far, CLI parsing worked just fine; no
		// need to show usage for runtime errors
		cmd.SilenceUsage = true

		return nil
	},
}

func Execute(ver string, gitsha string) int {
	rootCmd.SetVersionTemplate(`{{.Name}}/{{.Version}}`)
	rootCmd.Version = fmt.Sprintf("%s os/%s arch/%s",
		ver,
		runtime.GOOS,
		runtime.GOARCH,
	)

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetIn(os.Stdin)

	// rootCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
	// 	return err
	// })

	// hide the help flag as it's ubiquitous and thus noisy
	// we'll add it in the last line of the usage template
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	rootCmd.PersistentFlags().Lookup("help").Hidden = true

	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)

	cmd, err := rootCmd.ExecuteContextC(context.TODO())
	if err != nil {
		stop()

		// this is for the env command
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		if cmdutils.IsUserCancellation(err) {
			fmt.Fprintln(rootCmd.OutOrStdout(), "Exiting..")
			return 0
		}

		cmd.PrintErrln(utils.ErrorStyle.Render(cmd.ErrPrefix(), err.Error()))

		if cmdutils.IsNonUsageError(err) || cmdutils.IsAWSError(err) {
			return 1
		}

		if !cmd.Root().SilenceUsage {
			cmd.PrintErrf("Run '%v --help' for usage.\n", cmd.CommandPath())
		}
		// os.Exit(1)
		return 1
	}

	return 0
}
