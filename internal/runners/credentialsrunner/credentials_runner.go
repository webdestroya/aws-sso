package credentialsrunner

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/runners/loginrunner"
	"github.com/webdestroya/awssso/internal/utils"
	"github.com/webdestroya/awssso/internal/utils/awsutils"
)

func RunE(cmd *cobra.Command, args []string) error {

	profile := args[0]

	credentials, err := GetAWSCredentials(cmd.Context(), cmd.OutOrStdout(), profile)
	if err != nil {
		return err
	}

	creds := map[string]any{
		"Version":         1,
		"AccessKeyId":     credentials.AccessKeyID,
		"SecretAccessKey": credentials.SecretAccessKey,
		"SessionToken":    credentials.SessionToken,
		"Expiration":      credentials.Expires.Format(time.RFC3339),
	}

	pretty, err := utils.JsonifyPretty(creds)
	if err != nil {
		return err
	}

	cmd.Println(pretty)

	return nil
}

// func GetAWSCredentials(ctx context.Context, out io.Writer, profile string) (*ssoTypes.RoleCredentials, error) {
func GetAWSCredentials(ctx context.Context, out io.Writer, profile string) (*aws.Credentials, error) {

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

		tokenInfo, err = loginrunner.DoLoginFlow(ctx, out, ssoSession)
		if err != nil {
			return nil, err
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ssoSession.SSORegion))
	if err != nil {
		return nil, err
	}

	ssoClient := sso.NewFromConfig(cfg)
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

	// var provider aws.CredentialsProvider
	// provider = ssocreds.New(ssoClient, cfg.SSOAccountID, cfg.SSORoleName, cfg.SSOSession.SSOStartURL, func(options *ssocreds.Options) {
	// 	options.SSOTokenProvider = ssocreds.NewSSOTokenProvider(ssoOidcClient, tokenPath, func(o *ssocreds.SSOTokenProviderOptions) {
	// 		o.ClientOptions = append(o.ClientOptions, func(o2 *ssooidc.Options) {
	// 			o2.Region = ssoRegion
	// 			o2.Logger = logImpl
	// 		})
	// 	})
	// })

	// provider = aws.NewCredentialsCache(provider)

	// credentials, err := provider.Retrieve(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// return &credentials, nil

}
