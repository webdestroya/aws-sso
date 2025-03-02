package loginrunner

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/factory"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
)

type LoginOptions struct {
	NoBrowser bool
}

// https://github.com/synfinatic/aws-sso-cli/blob/main/internal/sso/awssso.go
// https://github.com/aws/aws-cli/blob/v2/awscli/customizations/sso/utils.py

func NewLoginCmd(f *factory.Factory) *cobra.Command {

	opts := &LoginOptions{}

	cmd := &cobra.Command{
		Use:               "login [PROFILE...]",
		Short:             "Login to the SSO session for the specified profile(s)",
		SilenceUsage:      true,
		ValidArgsFunction: profilepicker.ProfileCompletions,
		Args:              profilepicker.ValidProfileArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunE(opts, cmd, args)
		},
	}

	cmd.Flags().BoolVar(&opts.NoBrowser, "no-browser", false, "Disable opening the browser automatically")

	return cmd
}
