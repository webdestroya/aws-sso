package cmdutils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
)

func NewSSOOIDCClient(cfg aws.Config, optFns ...func(*ssooidc.Options)) *ssooidc.Client {
	return ssooidc.NewFromConfig(cfg, optFns...)
}

func NewSSOClient(cfg aws.Config, optFns ...func(*sso.Options)) *sso.Client {
	return sso.NewFromConfig(cfg, optFns...)
}
