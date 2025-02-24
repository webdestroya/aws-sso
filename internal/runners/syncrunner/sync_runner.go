package syncrunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/getcreds"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
	"github.com/webdestroya/aws-sso/internal/utils"
	"gopkg.in/ini.v1"
)

const (
	keyAccessKey    = `aws_access_key_id`
	keySecretKey    = `aws_secret_access_key`
	keySessionToken = `aws_session_token`
)

func RunE(opts *SyncOptions, cmd *cobra.Command, args []string) error {

	profiles, err := profilepicker.GetProfilesFromArgsOrPrompt(cmd, args)
	if err != nil {
		return err
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

	for _, profile := range profiles {
		if err := opts.syncCredentials(cmd.Context(), cmd.OutOrStdout(), credsIni, profile); err != nil {
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

func (opts *SyncOptions) syncCredentials(ctx context.Context, out io.Writer, credsIni *ini.File, profile string) error {

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

	creds, err := getcreds.GetAWSCredentials(ctx, out, profile)
	if err != nil {
		return err
	}

	credsIni.DeleteSection(profile)

	newSect, err := credsIni.NewSection(profile)
	if err != nil {
		return err
	}

	ak, _ := newSect.NewKey(keyAccessKey, creds.AccessKeyID)
	if creds.CanExpire {
		// time.Parse()
		ak.Comment = fmt.Sprintf("Expires: %s", creds.Expires.Format(time.RFC3339))
	}
	newSect.NewKey(keySecretKey, creds.SecretAccessKey)
	newSect.NewKey(keySessionToken, creds.SessionToken)

	return nil
}
