package syncrunner

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/webdestroya/aws-sso/internal/runners/credentialsrunner"
	"github.com/webdestroya/aws-sso/internal/utils"
	"gopkg.in/ini.v1"
)

const (
	keyAccessKey    = `aws_access_key_id`
	keySecretKey    = `aws_secret_access_key`
	keySessionToken = `aws_session_token`
	// keyRegion       = `region`
)

func RunE(cmd *cobra.Command, args []string) error {

	iniOpts := ini.LoadOptions{
		Loose:             true,
		AllowNestedValues: true,
	}

	credsFile := viper.GetString("sync.credentials_path")

	credsIni, err := ini.LoadSources(iniOpts, credsFile)
	if err != nil {
		return err
	}

	viper.BindPFlag("sync.force", cmd.Flag("force"))

	for _, profile := range args {
		if err := syncCredentials(cmd.Context(), cmd.OutOrStdout(), credsIni, profile); err != nil {
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

func syncCredentials(ctx context.Context, out io.Writer, credsIni *ini.File, profile string) error {

	// region := ""

	if credsIni.HasSection(profile) {
		sect, err := credsIni.GetSection(profile)
		if err != nil {
			return err
		}

		// TODO: allow this to be ignored
		// check to make sure the existing profile isnt something else
		if !viper.GetBool("sync.force") {
			if !(sect.HasKey(keyAccessKey) && sect.HasKey(keySecretKey) && sect.HasKey(keySessionToken)) {
				return fmt.Errorf("Profile %s already exists, but does not have AccessKey/SecretAccesKey/SessionToken. It probably is not an SSO profile.", profile)
			}
		}

		// if v, err := sect.GetKey(keyRegion); err == nil {
		// 	region = v.MustString(region)
		// }

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

	newSect.NewKey(keyAccessKey, creds.AccessKeyID)
	newSect.NewKey(keySecretKey, creds.SecretAccessKey)
	newSect.NewKey(keySessionToken, creds.SessionToken)
	// if region != "" {
	// 	newSect.NewKey(keyRegion, region)
	// }

	return nil
}
