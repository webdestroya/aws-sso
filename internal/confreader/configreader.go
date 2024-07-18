package confreader

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
)

func ReadAwsConfig() {
	fmt.Printf("FILES: %v\n", config.DefaultSharedConfigFiles)
	fmt.Printf("CREDS: %v\n", config.DefaultSharedCredentialsFiles)
}
