package listrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
)

type listOptions struct {
	ShowSessions bool
}

func NewListCmd(f *factory.Factory) *cobra.Command {

	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:               "list",
		Short:             "Lists available SSO profiles and any info about them",
		Args:              cobra.MatchAll(cobra.NoArgs),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			pl := &profileLister{
				cmd:  cmd,
				opts: opts,
			}

			return pl.run()
		},
	}

	cmd.Flags().BoolVar(&opts.ShowSessions, "sessions", false, "Show [sso-session *] entries")

	return cmd
}
