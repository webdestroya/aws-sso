package loginflow

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
	"github.com/toqueteos/webbrowser"
	"github.com/webdestroya/aws-sso/internal/utils"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

var (
	ErrAuthTimeoutError = errors.New("timeout waiting for authorization")
)

const (
	grantType  = `urn:ietf:params:oauth:grant-type:device_code`
	clientType = `public`
	clientName = `webdestroya-aws-sso-go`
)

type loginFlowOptions struct {
	DisableBrowser bool
}

type LoginFlowOption func(*loginFlowOptions)

func WithBrowser(v bool) LoginFlowOption {
	return func(lfo *loginFlowOptions) {
		lfo.DisableBrowser = !v
	}
}

func DoLoginFlow(ctx context.Context, out io.Writer, session *config.SSOSession, optFns ...LoginFlowOption) (*awsutils.AwsTokenInfo, error) {

	options := &loginFlowOptions{}
	for _, optFn := range optFns {
		optFn(options)
	}

	tokenFile, err := awsutils.GetSSOCachePath(session)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(session.SSORegion))
	if err != nil {
		return nil, err
	}

	client := cmdutils.NewSSOOIDCClient(cfg)
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

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Logging in to: %s\n", utils.CoalesceString(session.Name, session.SSOStartURL))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Attempting to automatically open the SSO authorization page in your default browser.")
	fmt.Fprintln(out, "If the browser does not open or you wish to use a different device to authorize this request,")
	fmt.Fprintln(out, "open the following URL:")
	fmt.Fprintln(out)
	fmt.Fprintln(out, verifUrl)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Then enter the code:")
	fmt.Fprintln(out)
	// https://github.com/common-nighthawk/go-figure
	fmt.Fprintf(out, "  %s\n", *sdaResp.UserCode)
	fmt.Fprintln(out)
	// fmt.Fprintf(out, "Successully logged into Start URL: %s\n", session.SSOStartURL)
	// fmt.Fprintln(out)

	if !options.DisableBrowser {
		go func() {
			realUrl := utils.CoalesceString(*sdaResp.VerificationUriComplete, *sdaResp.VerificationUri, verifUrl)
			if berr := webbrowser.Open(realUrl); berr != nil {
				fmt.Fprintln(out)

				fmt.Fprintln(out, utils.WarningStyle.Render(fmt.Sprintf("WARNING: Failed to open browser URL: %v", berr.Error())))
				fmt.Fprintln(out, "You will need to visit the verification URL and enter the code manually.")
				fmt.Fprintln(out)
			}
		}()
	}

	createTokenInput := &ssooidc.CreateTokenInput{
		ClientId:     regResp.ClientId,
		ClientSecret: regResp.ClientSecret,
		GrantType:    aws.String(grantType),
		DeviceCode:   sdaResp.DeviceCode,
	}

	var tokenObj *awsutils.AwsTokenInfo = nil

	accessTokenInvalidAfter := time.Now().Add(time.Duration(sdaResp.ExpiresIn) * time.Second)

	// spinCtx, cancelSpin := context.WithCancel(ctx)
	// defer cancelSpin()

	// spin := spinner.New().Context(spinCtx).Title("Waiting for authentication...")

	// go func() {
	// 	e := spin.Run()
	// 	fmt.Printf("SPIN ERR=%T msg=%v\n", e, e)
	// }()

	fmt.Fprintln(out, "Waiting for authentication...")

	for {
		// token expired dont bother
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

		tokenObj = &awsutils.AwsTokenInfo{
			AccessToken:  *createTokenOut.AccessToken,
			ExpiresAt:    &expiresAt,
			RefreshToken: aws.ToString(createTokenOut.RefreshToken),
			ClientID:     *regResp.ClientId,
			ClientSecret: *regResp.ClientSecret,
		}
		break
	}

	// cancelSpin()

	fmt.Fprintln(out, "Authentication successful! Writing token...")
	fmt.Fprintln(out)

	if err := awsutils.WriteAWSToken(tokenFile, *tokenObj); err != nil {
		return nil, err
	}

	return tokenObj, nil
}

// https://github.com/boto/botocore/blob/v2/botocore/utils.py
// https://github.com/mrtc0/aws-sso-go/blob/master/main.go
// https://github.com/boto/botocore/blob/v2/botocore/credentials.py#L2008
