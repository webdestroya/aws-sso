package awsutils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"gopkg.in/ini.v1"
)

func AddCredentialsToIni(file *ini.File, profile string, region string, creds aws.Credentials) error {

	file.DeleteSection(profile)

	newSect, err := file.NewSection(profile)
	if err != nil {
		return err
	}

	newSect.NewKey("region", region)
	newSect.NewKey("aws_access_key_id", creds.AccessKeyID)
	newSect.NewKey("aws_secret_access_key", creds.SecretAccessKey)
	newSect.NewKey("aws_session_token", creds.SessionToken)

	return nil
}
