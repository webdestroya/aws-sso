package awsutils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type IConfigLoader interface {
	LoadSharedConfigProfile(context.Context, string) (config.SharedConfig, error)
	LoadDefaultConfig(context.Context, ...func(*config.LoadOptions) error) (aws.Config, error)
}

var ConfigLoader IConfigLoader = defaultConfigLoader{}

func LoadSharedConfigProfile(ctx context.Context, profileName string) (config.SharedConfig, error) {
	return ConfigLoader.LoadSharedConfigProfile(ctx, profileName)
}

func LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return ConfigLoader.LoadDefaultConfig(ctx, optFns...)
}

type defaultConfigLoader struct{}

func (defaultConfigLoader) LoadSharedConfigProfile(ctx context.Context, profileName string) (config.SharedConfig, error) {
	return config.LoadSharedConfigProfile(ctx, profileName, func(lsco *config.LoadSharedConfigOptions) {
		// lsco.ConfigFiles = []string{}
		lsco.Logger = NewLogNone()
	})
}

func (defaultConfigLoader) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, optFns...)
}
