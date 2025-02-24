package listrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
)

func NewListCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "Lists available SSO sessions and any info about them",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.MatchAll(cobra.NoArgs),
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE:                  RunE,
	}

	return cmd
}
