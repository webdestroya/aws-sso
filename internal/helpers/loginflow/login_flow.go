package loginflow

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/webdestroya/aws-sso/internal/utils/awsutils"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

var (
	ErrAuthTimeoutError = cmdutils.NewNonUsageError("timeout waiting for authorization")
	ErrAuthDeniedError  = cmdutils.NewNonUsageError("authorization was denied")
)

const (
	grantType  = `urn:ietf:params:oauth:grant-type:device_code`
	clientType = `public`
	clientName = `webdestroya-aws-sso-go`
)

type loginFlowOptions struct {
	DisableBrowser bool
	UseDeviceCode  bool
}

type LoginFlowOption func(*loginFlowOptions)

func WithBrowser(v bool) LoginFlowOption {
	return func(lfo *loginFlowOptions) {
		lfo.DisableBrowser = !v
	}
}

func WithDeviceCode(v bool) LoginFlowOption {
	return func(lfo *loginFlowOptions) {
		lfo.UseDeviceCode = v
	}
}

func DoLoginFlow(ctx context.Context, out io.Writer, session *config.SSOSession, optFns ...LoginFlowOption) (*awsutils.AwsTokenInfo, error) {

	options := &loginFlowOptions{}
	for _, optFn := range optFns {
		optFn(options)
	}

	if options.UseDeviceCode {
		return doLoginFlowDeviceCode(ctx, out, session, options)
	}

	return doLoginFlowPKCE(ctx, out, session, options)

}

// https://github.com/boto/botocore/blob/v2/botocore/utils.py
// https://github.com/mrtc0/aws-sso-go/blob/master/main.go
// https://github.com/boto/botocore/blob/v2/botocore/credentials.py#L2008
