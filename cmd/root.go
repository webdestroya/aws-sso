package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/webdestroya/aws-sso/internal/appconfig"
	"github.com/webdestroya/aws-sso/internal/runners/rootrunner"
)

var (
	cfgFile        string
	verboseLogging = false
)

var rootCmd = &cobra.Command{
	Use:   "awssso",
	Short: "Facilitates usage of AWS sso authentication for older apps",
	RunE:  rootrunner.RunE,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func Execute(ver string, gitsha string) {
	rootCmd.SetVersionTemplate(`{{.Version}}`)
	rootCmd.Version = fmt.Sprintf("awssso/%s go/%s os/%s arch/%s",
		ver,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	err := rootCmd.ExecuteContext(context.TODO())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.awssso.toml)")
	rootCmd.MarkPersistentFlagFilename("config", "toml")
	rootCmd.PersistentFlags().MarkHidden("config")

	// rootCmd.PersistentFlags().BoolVar(&verboseLogging, "verbose", false, "Verbose logging")
	rootCmd.Flags().Bool("login", false, "Automatically login to profile")
	rootCmd.Flags().MarkHidden("login")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".aws"))
		viper.SetConfigType("toml")
		viper.SetConfigName(".awssso")
	}

	appconfig.SetDefaults(viper.GetViper())

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			fmt.Fprintf(os.Stderr, "Failed to load config: %s\n\n", err.Error())
		}
	}
}
