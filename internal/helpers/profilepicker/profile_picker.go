package profilepicker

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/listpicker"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

func PickSingleProfile(cmd *cobra.Command) (string, error) {
	profiles := Profiles()
	if len(profiles) == 0 {
		return "", cmdutils.NewNonUsageError("No SSO profile configurations found")
	}
	return listpicker.NewSingleChoice("Please select an AWS config profile:", profiles)
}

func GetProfilesFromArgsOrPrompt(cmd *cobra.Command, args []string) ([]string, error) {

	if len(args) == 0 {
		rep, err := PickSingleProfile(cmd)
		// cmd.Printf("RESULT: [%v] [err=%T/%v/]\n", rep, err, err)
		// if err != nil {
		// 	if errors.Is(err, huh.ErrUserAborted) {
		// 		os.Exit(0)
		// 		return nil, nil
		// 	}
		// }
		return []string{rep}, err
	}

	return args, nil
}
