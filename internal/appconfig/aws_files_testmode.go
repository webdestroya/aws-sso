//go:build testmode

package appconfig

func SetAwsCredentialPath(v string) {
	testModeCredentialsPath = v
}
