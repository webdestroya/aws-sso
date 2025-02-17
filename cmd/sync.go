package cmd

import (
	"errors"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/awslogger"
	"github.com/webdestroya/awssso/internal/awsutils"
	"gopkg.in/ini.v1"
)

var syncCmd = &cobra.Command{
	Use:          "sync [profile...]",
	Short:        "Sync AWS credentials. (This will overwrite the profile credentials!)",
	Example:      "awssso sync mycompany-production",
	Args:         cobra.ArbitraryArgs,
	RunE:         runSyncCmd,
	SilenceUsage: true,
}

var (
	errSyncsFailed   = errors.New("Some profiles failed to sync")
	errNotSSOProfile = errors.New("Profile is not an SSO-based profile")
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSyncCmd(cmd *cobra.Command, args []string) error {

	profilesList := slices.Compact(args)

	failCount := 0

	for _, profileName := range profilesList {

		cmd.Println()
		cmd.Printf("Syncing Profile: %s...\n", profileName)
		err := syncProfile(cmd, profileName)
		if err != nil {
			failCount += 1

			// cmd.Println(errorStatus)
			// cmd.Printf("  Error: %s\n", err.Error())
			cmd.Printf("  ")
			cmd.Println(errorStyle.Render("ERROR:", err.Error()))

		} else {
			// cmd.Println(statusSuccess)
			cmd.Printf("  ")
			cmd.Println(successStyle.Render("SUCCESS!"))
		}
	}

	// awsclient.LoadStuff()

	if failCount > 0 {
		return errSyncsFailed
	}

	return nil
}

func syncProfile(cmd *cobra.Command, profile string) error {

	// cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	// if err != nil {
	// 	return err
	// }

	// cmd.Printf("Profile: %s\n", profile)
	// cmd.Printf("CredType: %T\n", cfg.Credentials)
	// cmd.Printf("Type: %T\n", cfg.BearerAuthTokenProvider)
	// // cmd.Printf("ConfigSources: %v\n", cfg.ConfigSources...)

	// cmd.Println("Config Sources:")
	// for ci, confSource := range cfg.ConfigSources {
	// 	cmd.Printf("  %d: %T\n", ci, confSource)
	// }

	// cmd.Println()

	cfg, err := awsutils.LoadSharedConfigProfile(cmd.Context(), profile)
	if err != nil {
		return err
	}

	if cfg.SSOSession == nil {
		// cmd.Printf(errorStyle.Render(fmt.Sprintf("ERROR: Profile '%s' is not an SSO-based AWS profile.", cfg.Profile)))
		return errNotSSOProfile
	}

	ssoRegion := cfg.SSOSession.SSORegion

	vPrintf(cmd, "  Profile:  %s\n", cfg.Profile)
	vPrintf(cmd, "  SSO Name: %s\n", cfg.SSOSessionName)
	vPrintf(cmd, "  Region:   %s\n", ssoRegion)

	defConfig, err := config.LoadDefaultConfig(cmd.Context(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(ssoRegion),
		config.WithLogger(awslogger.NewLogNone()),
	)
	if err != nil {
		return err
	}

	iniOpts := ini.LoadOptions{
		Loose:             true,
		AllowNestedValues: true,
	}

	credsIni, err := ini.LoadSources(iniOpts, config.DefaultSharedCredentialsFilename())
	if err != nil {
		return err
	}

	ssoClient := sso.NewFromConfig(defConfig)
	ssoOidcClient := ssooidc.NewFromConfig(defConfig)
	tokenPath, err := ssocreds.StandardCachedTokenFilepath(cfg.SSOSessionName)
	if err != nil {
		return err
	}

	var provider aws.CredentialsProvider
	provider = ssocreds.New(ssoClient, cfg.SSOAccountID, cfg.SSORoleName, cfg.SSOSession.SSOStartURL, func(options *ssocreds.Options) {
		options.SSOTokenProvider = ssocreds.NewSSOTokenProvider(ssoOidcClient, tokenPath, func(o *ssocreds.SSOTokenProviderOptions) {
			o.ClientOptions = append(o.ClientOptions, func(o2 *ssooidc.Options) {
				o2.Region = ssoRegion
				o2.Logger = awslogger.NewLogNone()
			})
		})
	})

	// provider = aws.NewCredentialsCache(provider)

	credsIni.DeleteSection(profile)

	newSect, err := credsIni.NewSection(profile)
	if err != nil {
		return err
	}

	creds, err := provider.Retrieve(cmd.Context())
	if err != nil {
		return err
	}

	newSect.NewKey("aws_access_key_id", creds.AccessKeyID)
	newSect.NewKey("aws_secret_access_key", creds.SecretAccessKey)
	newSect.NewKey("aws_session_token", creds.SessionToken)

	credsIni.SaveTo("dummy.ini")

	return nil
}

func withLogger(lsco *config.LoadSharedConfigOptions) {
	lsco.Logger = nil
}
