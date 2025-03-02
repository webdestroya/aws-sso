package awsutils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadSharedConfigProfile(ctx context.Context, profileName string) (config.SharedConfig, error) {
	return config.LoadSharedConfigProfile(ctx, profileName, func(lsco *config.LoadSharedConfigOptions) {
		lsco.Logger = NewLogNone()
	})
}
