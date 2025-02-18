package awsutils

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadSharedConfigProfile(ctx context.Context, profileName string) (config.SharedConfig, error) {
	return config.LoadSharedConfigProfile(ctx, profileName, func(lsco *config.LoadSharedConfigOptions) {
		lsco.Logger = NewLogNone()
	})
}

func GetCredentialPath() string {
	if v, ok := os.LookupEnv("AWS_CREDENTIALS_PATH"); ok {
		return v
	}
	return config.DefaultSharedCredentialsFilename()
}
