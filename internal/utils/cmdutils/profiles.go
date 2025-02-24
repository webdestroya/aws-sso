package cmdutils

import (
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/profilepicker"
)

func GetProfilesFromArgsOrPrompt(cmd *cobra.Command, args []string) ([]string, error) {

	if len(args) == 0 {
		rep, err := profilepicker.PickSingleProfile(cmd)
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
