package appconfig

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

func SetDefaults(v *viper.Viper) *viper.Viper {
	viper.SetDefault("sync.force", false)
	viper.SetDefault("sync.credentials_path", config.DefaultSharedCredentialsFilename())
	return v
}
