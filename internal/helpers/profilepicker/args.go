package profilepicker

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var (
	ErrNoProfileProvidedError = errors.New("no profile argument was provided")
)

func ValidSingleProfileArg(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return ErrNoProfileProvidedError
	}
	if len(args) > 1 {
		return fmt.Errorf("only one profile can be provided")
	}

	profiles := Profiles()

	if len(profiles) == 0 {

	}

	profile := args[0]
	if !slices.Contains(profiles, profile) {
		return fmt.Errorf("invalid argument: %s is not an SSO profile", profile)
	}

	return nil
}

func ValidProfileArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}

	profiles := Profiles()

	for _, profile := range args {
		if !slices.Contains(profiles, profile) {
			return fmt.Errorf("invalid argument: %s is not an SSO profile", profile)
		}
	}
	return nil
}

func ValidProfileFirstArg(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return ErrNoProfileProvidedError
	}

	profile := args[0]
	if !slices.Contains(Profiles(), profile) {
		return fmt.Errorf("invalid argument: %s is not an SSO profile", profile)
	}
	return nil
}
