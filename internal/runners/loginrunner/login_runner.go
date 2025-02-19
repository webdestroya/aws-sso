// This will perform the login function to a specific StartURL SSO endpoint
// it will then update the cached token for the StartURL and/or Name
// it then exits, it does not actually get role credentials
package loginrunner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	ssooidcTypes "github.com/aws/aws-sdk-go-v2/service/ssooidc/types"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
)

var (
	ErrAuthTimeoutError = errors.New("timeout waiting for authorization")
)

type Printer interface {
	Printf(string, ...any)
	Println(...any)
}

const (
	grantType  = `urn:ietf:params:oauth:grant-type:device_code`
	clientType = `public`
	clientName = `webdestroya-aws-sso-go`
)

// type Runner struct {
// 	cmd     *cobra.Command
// 	cfg     *config.SharedConfig
// 	defCfg  *aws.Config
// 	profile string
// }

func RunE(cmd *cobra.Command, args []string) error {

	// TODO: iterate all the profiles and make sure they are actually SSO things
	// TODO: reduce to a unique list of start_urls
	// login to each one

	ctx := cmd.Context()

	cfgmap := make(map[string]*config.SSOSession)

	for _, profile := range args {
		sharedCfg, err := awsutils.LoadSharedConfigProfile(ctx, profile)
		if err != nil {
			return err
		}

		ssoSession, err := awsutils.ExtractSSOInfo(sharedCfg)
		if err != nil {
			return err
		}

		cachePath, err := awsutils.GetSSOCachePath(ssoSession)
		if err != nil {
			return err
		}

		cfgmap[cachePath] = ssoSession
	}

	for _, session := range cfgmap {
		if _, err := DoLoginFlow(ctx, cmd.OutOrStdout(), session); err != nil {
			return err
		}
	}

	return nil
}

func DoLoginFlow(ctx context.Context, out io.Writer, session *config.SSOSession) (*awsutils.AwsTokenInfo, error) {
	tokenFile, err := awsutils.GetSSOCachePath(session)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(session.SSORegion))
	if err != nil {
		return nil, err
	}

	client := ssooidc.NewFromConfig(cfg)
	regResp, err := client.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String(clientName),
		ClientType: aws.String(clientType),
		GrantTypes: []string{grantType},
	})
	if err != nil {
		return nil, err
	}

	sdaResp, err := client.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     regResp.ClientId,
		ClientSecret: regResp.ClientSecret,
		StartUrl:     &session.SSOStartURL,
	})
	if err != nil {
		return nil, err
	}

	verifUrl := fmt.Sprintf("https://device.sso.%s.amazonaws.com/", cfg.Region)

	fmt.Fprintf(out, "Logging in to: %s\n", utils.CoalesceString(session.Name, session.SSOStartURL))
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Attempting to automatically open the SSO authorization page in your default browser.")
	fmt.Fprintln(out, "If the browser does not open or you wish to use a different device to authorize this request,")
	fmt.Fprintln(out, "open the following URL:")
	fmt.Fprintln(out)
	fmt.Fprintln(out, verifUrl)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Then enter the code:")
	fmt.Fprintln(out)
	fmt.Fprintln(out, *sdaResp.UserCode)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Successully logged into Start URL: %s\n", session.SSOStartURL)
	fmt.Fprintln(out)

	go func() {
		realUrl := utils.CoalesceString(*sdaResp.VerificationUriComplete, *sdaResp.VerificationUri, verifUrl)
		_ = webbrowser.Open(realUrl)
	}()

	createTokenInput := &ssooidc.CreateTokenInput{
		ClientId:     regResp.ClientId,
		ClientSecret: regResp.ClientSecret,
		GrantType:    aws.String(grantType),
		DeviceCode:   sdaResp.DeviceCode,
	}

	var tokenObj *awsutils.AwsTokenInfo = nil

	accessTokenInvalidAfter := time.Now().Add(time.Duration(sdaResp.ExpiresIn) * time.Second)

	for {
		if time.Now().After(accessTokenInvalidAfter) {
			return nil, ErrAuthTimeoutError
		}

		createTokenOut, err := client.CreateToken(ctx, createTokenInput)
		if err != nil {
			var authPendingError *ssooidcTypes.AuthorizationPendingException
			if errors.As(err, &authPendingError) {
				time.Sleep(time.Duration(sdaResp.Interval) * time.Second)
				continue
			} else {
				return nil, err
			}
		}

		expiresAt := awsutils.RFC3339(time.Now().Add(time.Duration(createTokenOut.ExpiresIn) * time.Second))

		// accessToken = *createTokenOut.AccessToken
		tokenObj = &awsutils.AwsTokenInfo{
			AccessToken:  *createTokenOut.AccessToken,
			ExpiresAt:    &expiresAt,
			RefreshToken: aws.ToString(createTokenOut.RefreshToken),
			ClientID:     *regResp.ClientId,
			ClientSecret: *regResp.ClientSecret,
		}
		break
	}

	// if tokenObj == nil {
	// 	return nil, errors.New("could not get token")
	// }

	if err := awsutils.WriteAWSToken(tokenFile, *tokenObj); err != nil {
		return nil, err
	}

	return tokenObj, nil
}

// https://github.com/boto/botocore/blob/v2/botocore/utils.py
// https://github.com/mrtc0/aws-sso-go/blob/master/main.go
/*
func (r *Runner) DoLoginFlow(ctx context.Context) (*ssoTypes.RoleCredentials, error) {

	sharedCfg, err := awsutils.LoadSharedConfigProfile(ctx, r.profile)
	if err != nil {
		return nil, err
	}

	ssoSession, err := awsutils.ExtractSSOInfo(sharedCfg)
	if err != nil {
		return nil, err
	}

	tokenFile, err := ssocreds.StandardCachedTokenFilepath(ssoSession.SSOStartURL)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ssoSession.SSORegion))
	if err != nil {
		return nil, err
	}

	client := ssooidc.NewFromConfig(cfg)
	regResp, err := client.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String("acorns-aws-sso-go"),
		ClientType: aws.String(clientType),
		GrantTypes: []string{grantType},
		// EntitledApplicationArn: new(string),
		// IssuerUrl:              new(string),
		// RedirectUris:           []string{},
		// Scopes:                 []string{},
	})
	if err != nil {
		return nil, err
	}

	sdaResp, err := client.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     regResp.ClientId,
		ClientSecret: regResp.ClientSecret,
		StartUrl:     &ssoSession.SSOStartURL,
	})
	if err != nil {
		return nil, err
	}

	verifUrl := fmt.Sprintf("https://device.sso.%s.amazonaws.com/", cfg.Region)

	r.cmd.Println("Attempting to automatically open the SSO authorization page in your default browser.")
	r.cmd.Println("If the browser does not open or you wish to use a different device to authorize this request,")
	r.cmd.Println("open the following URL:")
	r.cmd.Println()
	r.cmd.Println(verifUrl)
	r.cmd.Println()
	r.cmd.Println("Then enter the code:")
	r.cmd.Println()
	r.cmd.Println(*sdaResp.UserCode)
	r.cmd.Println()
	r.cmd.Printf("Successully logged into Start URL: %s\n", ssoSession.SSOStartURL)
	r.cmd.Println()

	go func() {
		realUrl := utils.CoalesceString(*sdaResp.VerificationUriComplete, *sdaResp.VerificationUri, verifUrl)
		_ = webbrowser.Open(realUrl)
	}()

	createTokenInput := &ssooidc.CreateTokenInput{
		ClientId:     regResp.ClientId,
		ClientSecret: regResp.ClientSecret,
		GrantType:    aws.String(grantType),
		DeviceCode:   sdaResp.DeviceCode,
	}

	var tokenObj *awsutils.AwsTokenInfo = nil

	accessTokenInvalidAfter := time.Now().Add(time.Duration(sdaResp.ExpiresIn) * time.Second)

	for {
		if time.Now().After(accessTokenInvalidAfter) {
			return nil, errors.New("timeout waiting for authorization")
		}

		createTokenOut, err := client.CreateToken(ctx, createTokenInput)
		if err != nil {
			var authPendingError *ssooidcTypes.AuthorizationPendingException
			if errors.As(err, &authPendingError) {
				time.Sleep(time.Duration(sdaResp.Interval) * time.Second)
				continue
			} else {
				return nil, err
			}
		}

		expiresAt := awsutils.RFC3339(time.Now().Add(time.Duration(createTokenOut.ExpiresIn) * time.Second))

		// accessToken = *createTokenOut.AccessToken
		tokenObj = &awsutils.AwsTokenInfo{
			AccessToken:  *createTokenOut.AccessToken,
			ExpiresAt:    &expiresAt,
			RefreshToken: aws.ToString(createTokenOut.RefreshToken),
			ClientID:     *regResp.ClientId,
			ClientSecret: *regResp.ClientSecret,
		}
		break
	}

	// if tokenObj == nil {
	// 	return nil, errors.New("could not get token")
	// }

	if err := awsutils.WriteAWSToken(tokenFile, *tokenObj); err != nil {
		return nil, err
	}

	if true {
		return nil, nil
	}

	ssoClient := sso.NewFromConfig(cfg)
	creds, err := ssoClient.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{
		AccessToken: &tokenObj.AccessToken,
		AccountId:   &sharedCfg.SSOAccountID,
		RoleName:    &sharedCfg.SSORoleName,
	})
	if err != nil {
		return nil, err
	}

	roleCreds := creds.RoleCredentials

	out, _ := json.Marshal(roleCreds)

	r.cmd.Println("JSON:", string(out))

	return roleCreds, nil
}
*/
