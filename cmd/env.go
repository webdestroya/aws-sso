package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:                   "env PROFILE -- command [command-args...]",
	Aliases:               []string{"run"},
	Short:                 "Run a command with AWS access keys set in the environment",
	SilenceUsage:          false,
	DisableFlagsInUseLine: true,
	FParseErrWhitelist:    cobra.FParseErrWhitelist{UnknownFlags: true},
	Args:                  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runEnv(cmd, args)
	},
	// RunE:                  runEnv,
	// DisableFlagParsing:    true,
}

func init() {
	rootCmd.AddCommand(envCmd)
}

func runEnv(cmd *cobra.Command, args []string) error {

	cmd.Printf("COMMAND: %v\n", args)

	return nil

	profile, args := args[0], args[1:]

	command, args := args[0], args[1:]

	if command == "--" {
		if len(args) < 2 {
			return errors.New("No command was provided")
		}
		command, args = args[0], args[1:]
	}

	binPath, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	credinfo, err := getAwsCreds(cmd, profile)
	if err != nil {
		return err
	}

	env := os.Environ()

	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", credinfo.Creds.AccessKeyID))
	env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", credinfo.Creds.SecretAccessKey))
	env = append(env, fmt.Sprintf("AWS_SESSION_TOKEN=%s", credinfo.Creds.SessionToken))
	env = append(env, fmt.Sprintf("AWS_REGION=%s", credinfo.Region))

	return syscall.Exec(binPath, args, env)
}
