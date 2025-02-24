package credentialsrunner

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/webdestroya/aws-sso/internal/runners/loginrunner"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

// func GetAWSCredentials(ctx context.Context, out io.Writer, profile string) (*ssoTypes.RoleCredentials, error) {
func GetAWSCredentials(ctx context.Context, out io.Writer, profile string, optFns ...loginrunner.LoginFlowOption) (*aws.Credentials, error) {

	sharedCfg, err := awsutils.LoadSharedConfigProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	ssoSession, err := awsutils.ExtractSSOInfo(sharedCfg)
	if err != nil {
		return nil, err
	}

	tokenFile, err := awsutils.GetSSOCachePath(ssoSession)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := awsutils.ReadTokenFile(tokenFile)
	if err != nil || tokenInfo == nil || tokenInfo.Expired() {

		// problem with getting token file, so do the login flow

		tokenInfo, err = loginrunner.DoLoginFlow(ctx, out, ssoSession, optFns...)
		if err != nil {
			return nil, err
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ssoSession.SSORegion))
	if err != nil {
		return nil, err
	}

	ssoClient := cmdutils.NewSSOClient(cfg)
	creds, err := ssoClient.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{
		AccessToken: &tokenInfo.AccessToken,
		AccountId:   &sharedCfg.SSOAccountID,
		RoleName:    &sharedCfg.SSORoleName,
	})
	if err != nil {
		return nil, err
	}

	roleCreds := creds.RoleCredentials
	return &aws.Credentials{
		AccessKeyID:     *roleCreds.AccessKeyId,
		SecretAccessKey: *roleCreds.SecretAccessKey,
		SessionToken:    *roleCreds.SessionToken,
		Source:          ssocreds.ProviderName,
		CanExpire:       true,
		Expires:         time.Unix(0, roleCreds.Expiration*int64(time.Millisecond)).UTC(),
		AccountID:       sharedCfg.SSOAccountID,
	}, nil
}
