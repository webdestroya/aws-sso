//go:build testmode

package appconfig

func GetAwsCredentialPath() string {
	return "test-credentials"
}

func GetAwsConfigPath() string {
	return "test-config"
}
