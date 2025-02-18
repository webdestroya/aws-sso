package syncrunner

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/runners/credentialsrunner"
	"gopkg.in/ini.v1"
)

func RunE(cmd *cobra.Command, args []string) error {

	iniOpts := ini.LoadOptions{
		Loose:             true,
		AllowNestedValues: true,
	}

	credsIni, err := ini.LoadSources(iniOpts, config.DefaultSharedCredentialsFilename())
	if err != nil {
		return err
	}

	for _, profile := range args {
		if err := syncCredentials(cmd.Context(), cmd.OutOrStdout(), credsIni, profile); err != nil {
			return err
		}
	}

	credsIni.WriteTo(cmd.OutOrStdout())

	return nil
}

// func SyncProfile()

func syncCredentials(ctx context.Context, out io.Writer, credsIni *ini.File, profile string) error {

	if credsIni.HasSection(profile) {
		sect, err := credsIni.GetSection(profile)
		if err != nil {
			return err
		}

		// TODO: allow this to be ignored
		// check to make sure the existing profile isnt something else
		if !(sect.HasKey("aws_access_key_id") && sect.HasKey("aws_secret_access_key") && sect.HasKey("aws_session_token")) {
			return fmt.Errorf("Profile %s already exists, but does not have AccessKey/SecretAccesKey/SessionToken. It probably is not an SSO profile.", profile)
		}

		credsIni.DeleteSection(profile)
	}

	newSect, err := credsIni.NewSection(profile)
	if err != nil {
		return err
	}

	creds, err := credentialsrunner.GetAWSCredentials(ctx, out, profile)
	if err != nil {
		return err
	}

	newSect.NewKey("aws_access_key_id", creds.AccessKeyID)
	newSect.NewKey("aws_secret_access_key", creds.SecretAccessKey)
	newSect.NewKey("aws_session_token", creds.SessionToken)

	return nil
}
