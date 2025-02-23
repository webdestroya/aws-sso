package syncrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/components/profilepicker"
	"github.com/webdestroya/aws-sso/internal/factory"
	"github.com/webdestroya/aws-sso/internal/helpers/profilelist"
)

type SyncOptions struct {
	Login                 bool
	Force                 bool
	CredentialsOutputPath string
}

func NewCmdSync(f *factory.Factory) *cobra.Command {

	opts := &SyncOptions{
		Login:                 true,
		Force:                 false,
		CredentialsOutputPath: appconfig.GetAwsCredentialPath(),
	}

	cmd := &cobra.Command{
		Use:               "sync [PROFILE...]",
		Short:             "Sync AWS credentials. (This will overwrite the profile credentials!)",
		Example:           "awssso sync mycompany-production",
		SilenceUsage:      true,
		ValidArgsFunction: profilelist.ProfileCompletions,
		RunE: func(c *cobra.Command, args []string) error {
			return RunE(opts, c, args)
		},
		PreRunE: profilepicker.EnsureProfileArgsPreRunE,
		// Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().BoolVar(&opts.Login, "login", true, "Automatically login to profile")
	cmd.Flags().MarkHidden("login")
	cmd.Flags().MarkDeprecated("login", "you will automatically be prompted if necessary")

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Force overwrite profile credentials even if they do not appear to be for an SSO profile")

	return cmd
}
