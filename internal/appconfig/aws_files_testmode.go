//go:build testmode

package appconfig

func SetAwsCredentialPath(v string) string {
	return testModeCredentialsPath = v
}
