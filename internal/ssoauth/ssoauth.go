package ssoauth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/smithy-go/logging"
	"github.com/webdestroya/awssso/internal/awslogger"
	"github.com/webdestroya/awssso/internal/util"
)

type SSOAuth struct {
	Ctx         context.Context
	ProfileName string

	ssoRegion string

	cfg       *config.SharedConfig
	defConfig aws.Config

	ssoClient     *sso.Client
	ssoOidcClient *ssooidc.Client
	tokenPath     string
	provider      aws.CredentialsProvider

	Logger logging.Logger
}

func New(ctx context.Context, profile string) (*SSOAuth, error) {

	logImpl := awslogger.NewLogNone()

	sa := &SSOAuth{
		Ctx:         ctx,
		ProfileName: profile,
		Logger:      logImpl,
	}

	return sa, nil
}

func (sa *SSOAuth) Init() error {
	cfg, err := config.LoadSharedConfigProfile(sa.Ctx, sa.ProfileName,
		func(lsco *config.LoadSharedConfigOptions) {
			lsco.Logger = sa.Logger
		},
	)
	if err != nil {
		return err
	}
	sa.cfg = &cfg

	sa.ssoRegion = util.CoalesceString(cfg.SSOSession.SSORegion, cfg.Region)

	sa.defConfig, err = config.LoadDefaultConfig(sa.Ctx,
		config.WithSharedConfigProfile(sa.ProfileName),
		config.WithRegion(sa.ssoRegion),
		config.WithLogger(sa.Logger),
	)
	if err != nil {
		return err
	}

	return nil
}

func (sa *SSOAuth) InitProvider() error {
	sa.ssoClient = sso.NewFromConfig(sa.defConfig)
	sa.ssoOidcClient = ssooidc.NewFromConfig(sa.defConfig)
	tokenPath, err := ssocreds.StandardCachedTokenFilepath(sa.cfg.SSOSessionName)
	if err != nil {
		return err
	}
	sa.tokenPath = tokenPath

	var provider aws.CredentialsProvider
	provider = ssocreds.New(sa.ssoClient, sa.cfg.SSOAccountID, sa.cfg.SSORoleName, sa.cfg.SSOSession.SSOStartURL, func(options *ssocreds.Options) {
		options.SSOTokenProvider = ssocreds.NewSSOTokenProvider(sa.ssoOidcClient, sa.tokenPath, func(o *ssocreds.SSOTokenProviderOptions) {
			o.ClientOptions = append(o.ClientOptions, func(o2 *ssooidc.Options) {
				o2.Region = sa.ssoRegion
				o2.Logger = sa.Logger
			})
		})
	})

	sa.provider = aws.NewCredentialsCache(provider)

	return nil
}

// https://github.com/aws/aws-cli/blob/v2/awscli/customizations/sso/utils.py
