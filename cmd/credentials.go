package cmd

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/awslogger"
	"github.com/webdestroya/awssso/internal/util"
)

// credentialsCmd represents the credentials command
var credentialsCmd = &cobra.Command{
	Use:          "credentials",
	Short:        "Use SSO creds as AWS process credentials",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         runCredentials,
}

func init() {
	rootCmd.AddCommand(credentialsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// credentialsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// credentialsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCredentials(cmd *cobra.Command, args []string) error {

	profileName := args[0]

	cfg, err := config.LoadSharedConfigProfile(cmd.Context(), profileName, withLogger)
	if err != nil {
		return err
	}

	ssoRegion := util.CoalesceString(cfg.SSOSession.SSORegion, cfg.Region)

	defConfig, err := config.LoadDefaultConfig(cmd.Context(),
		config.WithSharedConfigProfile(profileName),
		config.WithRegion(ssoRegion),
		config.WithLogger(awslogger.NewLogAll(cmd.OutOrStdout())),
	)
	if err != nil {
		return err
	}

	ssoClient := sso.NewFromConfig(defConfig)
	ssoOidcClient := ssooidc.NewFromConfig(defConfig)
	tokenPath, err := ssocreds.StandardCachedTokenFilepath(cfg.SSOSessionName)
	if err != nil {
		return err
	}

	cmd.Printf("DefRegion: %s\n", defConfig.Region)
	cmd.Printf("Acct:      %s\n", cfg.SSOAccountID)
	cmd.Printf("RoleName:  %s\n", cfg.SSORoleName)
	cmd.Printf("StartURL:  %s\n", cfg.SSOSession.SSOStartURL)
	cmd.Printf("Region:    %s\n", ssoRegion)

	var provider aws.CredentialsProvider
	provider = ssocreds.New(ssoClient, cfg.SSOAccountID, cfg.SSORoleName, cfg.SSOSession.SSOStartURL, func(options *ssocreds.Options) {
		options.SSOTokenProvider = ssocreds.NewSSOTokenProvider(ssoOidcClient, tokenPath, func(o *ssocreds.SSOTokenProviderOptions) {
			o.ClientOptions = append(o.ClientOptions, func(o2 *ssooidc.Options) {
				o2.Region = ssoRegion
				o2.Logger = awslogger.NewLogAll(cmd.OutOrStdout())
			})
		})
	})

	provider = aws.NewCredentialsCache(provider)

	credentials, err := provider.Retrieve(cmd.Context())
	if err != nil {
		return err
	}

	// https://docs.aws.amazon.com/sdkref/latest/guide/feature-process-credentials.html
	creds := map[string]any{
		"Version":         1,
		"AccessKeyId":     credentials.AccessKeyID,
		"SecretAccessKey": credentials.SecretAccessKey,
		"SessionToken":    credentials.SessionToken,
		"Expiration":      credentials.Expires.Format(time.RFC3339),
	}

	cmd.Println(util.JsonifyPretty(creds))

	return nil
}

type credInfo struct {
	Creds  *aws.Credentials
	Region string
}

func getAwsCreds(cmd *cobra.Command, profileName string) (*credInfo, error) {

	ctx := cmd.Context()
	logImpl := awslogger.NewLogAll(cmd.OutOrStdout())

	cfg, err := config.LoadSharedConfigProfile(ctx, profileName, withLogger)
	if err != nil {
		return nil, err
	}

	ssoRegion := util.CoalesceString(cfg.SSOSession.SSORegion, cfg.Region)

	defConfig, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profileName),
		config.WithRegion(ssoRegion),
		config.WithLogger(logImpl),
	)
	if err != nil {
		return nil, err
	}

	ssoClient := sso.NewFromConfig(defConfig)
	ssoOidcClient := ssooidc.NewFromConfig(defConfig)
	tokenPath, err := ssocreds.StandardCachedTokenFilepath(cfg.SSOSessionName)
	if err != nil {
		return nil, err
	}

	var provider aws.CredentialsProvider
	provider = ssocreds.New(ssoClient, cfg.SSOAccountID, cfg.SSORoleName, cfg.SSOSession.SSOStartURL, func(options *ssocreds.Options) {
		options.SSOTokenProvider = ssocreds.NewSSOTokenProvider(ssoOidcClient, tokenPath, func(o *ssocreds.SSOTokenProviderOptions) {
			o.ClientOptions = append(o.ClientOptions, func(o2 *ssooidc.Options) {
				o2.Region = ssoRegion
				o2.Logger = logImpl
			})
		})
	})

	provider = aws.NewCredentialsCache(provider)

	credentials, err := provider.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	return &credInfo{
		Creds:  &credentials,
		Region: ssoRegion,
	}, nil
}
