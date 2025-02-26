package getcreds

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/webdestroya/aws-sso/internal/helpers/loginflow"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

var (
	ErrTokenInvalidError = errors.New("token was not found or was expired. You need to login to this profile first.")
)

// func GetAWSCredentials(ctx context.Context, out io.Writer, profile string) (*ssoTypes.RoleCredentials, error) {
func GetAWSCredentials(ctx context.Context, out io.Writer, profile string, optFns ...GetCredOption) (*aws.Credentials, error) {

	opts := &getCredOptions{
		LoginFlowOptions: make([]loginflow.LoginFlowOption, 0),
	}
	for _, optFn := range optFns {
		optFn(opts)
	}

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

		if opts.DisableLogin {
			return nil, ErrTokenInvalidError
		}

		// problem with getting token file, so do the login flow

		tokenInfo, err = loginflow.DoLoginFlow(ctx, out, ssoSession, opts.LoginFlowOptions...)
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
	result := &aws.Credentials{
		AccessKeyID:     *roleCreds.AccessKeyId,
		SecretAccessKey: *roleCreds.SecretAccessKey,
		SessionToken:    *roleCreds.SessionToken,
		Source:          ssocreds.ProviderName,
		CanExpire:       true,
		Expires:         time.Unix(0, roleCreds.Expiration*int64(time.Millisecond)).UTC(),
		AccountID:       sharedCfg.SSOAccountID,
	}

	if opts.CliCache {
		if err := writeCliCache(sharedCfg, ssoSession, result); err != nil {
			return result, err
		}
	}

	return result, nil
}
