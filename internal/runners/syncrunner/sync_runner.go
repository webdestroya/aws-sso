package syncrunner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/runners/credentialsrunner"
	"github.com/webdestroya/aws-sso/internal/utils"
	"gopkg.in/ini.v1"
)

const (
	keyAccessKey    = `aws_access_key_id`
	keySecretKey    = `aws_secret_access_key`
	keySessionToken = `aws_session_token`
)

func RunE(opts *SyncOptions, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		return errors.New("No profiles were provided to sync")
	}

	iniOpts := ini.LoadOptions{
		Loose:             true,
		AllowNestedValues: true,
	}

	credsFile := opts.CredentialsOutputPath

	credsIni, err := ini.LoadSources(iniOpts, credsFile)
	if err != nil {
		return err
	}

	for _, profile := range args {
		if err := syncCredentials(opts, cmd.Context(), cmd.OutOrStdout(), credsIni, profile); err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)

	_, err = credsIni.WriteTo(buf)
	if err != nil {
		return err
	}

	return utils.AtomicWriteFile(credsFile, buf.Bytes(), 0600)
}

func syncCredentials(opts *SyncOptions, ctx context.Context, out io.Writer, credsIni *ini.File, profile string) error {

	if credsIni.HasSection(profile) {
		sect, err := credsIni.GetSection(profile)
		if err != nil {
			return err
		}

		// check to make sure the existing profile isnt something else
		if !opts.Force {
			if !(sect.HasKey(keyAccessKey) && sect.HasKey(keySecretKey) && sect.HasKey(keySessionToken)) {
				return fmt.Errorf("Profile %s already exists, but does not have AccessKey/SecretAccesKey/SessionToken. It probably is not an SSO profile.", profile)
			}
		}

	}

	creds, err := credentialsrunner.GetAWSCredentials(ctx, out, profile)
	if err != nil {
		return err
	}

	credsIni.DeleteSection(profile)

	newSect, err := credsIni.NewSection(profile)
	if err != nil {
		return err
	}

	newSect.NewKey(keyAccessKey, creds.AccessKeyID)
	newSect.NewKey(keySecretKey, creds.SecretAccessKey)
	newSect.NewKey(keySessionToken, creds.SessionToken)

	return nil
}
