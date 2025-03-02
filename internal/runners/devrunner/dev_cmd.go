//go:build !nodev

package devrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
)

func NewDevCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "dev",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		RunE:                  RunE,
	}

	return cmd
}
