package envrunner

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
)

type envOptions struct {
	WaitDelay time.Duration
	Quiet     bool
}

func NewEnvCmd(f *factory.Factory) *cobra.Command {

	opts := &envOptions{}

	cmd := &cobra.Command{
		Use:                   "env [flags] PROFILE -- command [command-args...]",
		Aliases:               []string{"run"},
		Short:                 "Run a command with AWS access keys set in the environment",
		DisableFlagsInUseLine: true,
		// FParseErrWhitelist:    cobra.FParseErrWhitelist{UnknownFlags: true},
		Args: cobra.MatchAll(profilepicker.ValidProfileFirstArg, cobra.MinimumNArgs(2)),
		RunE: func(c *cobra.Command, args []string) error {
			return opts.runE(c, args)
		},
	}

	cmd.Flags().SetInterspersed(false)
	cmd.Flags().DurationVar(&opts.WaitDelay, "wait-delay", 5*time.Second, "Max duration to wait after SIGTERM before sending SIGKILL")
	cmd.Flags().BoolVar(&opts.Quiet, "quiet", false, "Don't print info about the command to be run")

	// cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
	// 	c.Printf("ERR: %v\n", err.Error())
	// 	return nil
	// })

	return cmd
}
