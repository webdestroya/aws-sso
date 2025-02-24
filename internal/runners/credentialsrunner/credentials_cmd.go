package credentialsrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
)

func NewCredentialsCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "credentials PROFILE",
		Short:                 "Use SSO creds as AWS process credentials",
		DisableFlagsInUseLine: true,
		ValidArgsFunction:     profilepicker.ProfileCompletions,
		Args:                  cobra.MatchAll(profilepicker.ValidSingleProfileArg),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunE(cmd, args)
		},
		// SilenceUsage:          true,
	}

	return cmd
}
