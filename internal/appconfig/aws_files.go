//go:build !testmode

package appconfig

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAwsCredentialPath() string {
	// AWS_SHARED_CREDENTIALS_FILE
	// AWS_CREDENTIALS_PATH
	if v, ok := os.LookupEnv("AWS_SHARED_CREDENTIALS_FILE"); ok {
		return v
	}
	return config.DefaultSharedCredentialsFilename()
}

func GetAwsConfigPath() string {
	if v, ok := os.LookupEnv("AWS_CONFIG_FILE"); ok {
		return v
	}
	return config.DefaultSharedConfigFilename()
}
