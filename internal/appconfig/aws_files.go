package appconfig

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	testModeCredentialsPath = "internal/appconfig/testdata/credentials"
	testModeConfigPath      = "internal/appconfig/testdata/config"
)

func GetAwsCredentialPath() string {
	// AWS_SHARED_CREDENTIALS_FILE
	// AWS_CREDENTIALS_PATH
	if testing.Testing() {
		return testModeCredentialsPath
	}

	if v, ok := os.LookupEnv("AWS_SHARED_CREDENTIALS_FILE"); ok {
		return v
	}
	return config.DefaultSharedCredentialsFilename()
}

func GetAwsConfigPath() string {
	if testing.Testing() {
		return testModeConfigPath
	}

	if v, ok := os.LookupEnv("AWS_CONFIG_FILE"); ok {
		return v
	}
	return config.DefaultSharedConfigFilename()
}
