package cmd

import (
	"errors"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/confreader"
)

var syncCmd = &cobra.Command{
	Use:          "sync [profile...]",
	Short:        "Sync AWS credentials. (This will overwrite the profile credentials!)",
	Args:         cobra.ArbitraryArgs,
	RunE:         runSyncCmd,
	SilenceUsage: true,
}

var (
	errSyncsFailed   = errors.New("Some profiles failed to sync")
	errNotSSOProfile = errors.New("Profile is not an SSO-based profile")
)

var (
	errorStatus   = errorStyle.Render("ERROR")
	statusSuccess = successStyle.Render("SUCCESS")
)

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runSyncCmd(cmd *cobra.Command, args []string) error {

	slices.Sort(args)
	profilesList := slices.Compact(args)
	fmt.Fprintln(cmd.OutOrStdout(), "Heyooo", profilesList)

	confreader.ReadAwsConfig()

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

	cfg, err := config.LoadSharedConfigProfile(cmd.Context(), profile, withLogger)
	if err != nil {
		return err
	}

	vPrintf(cmd, "  Profile: %s\n", cfg.Profile)
	vPrintf(cmd, "  SSO Name: %s\n", cfg.SSOSessionName)
	vPrintf(cmd, "  SSO Session: %v\n", cfg.SSOSession)
	vPrintf(cmd, "  CredSource: %v\n", cfg.CredentialSource)
	vPrintf(cmd, "  Creds: %v\n", cfg.Credentials)

	if cfg.SSOSession == nil {
		// cmd.Printf(errorStyle.Render(fmt.Sprintf("ERROR: Profile '%s' is not an SSO-based AWS profile.", cfg.Profile)))
		return errNotSSOProfile
	}

	return nil
}

func withLogger(lsco *config.LoadSharedConfigOptions) {
	lsco.Logger = nil
}
