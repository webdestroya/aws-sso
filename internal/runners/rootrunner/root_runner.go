package rootrunner

import "github.com/spf13/cobra"

func RunE(cmd *cobra.Command, args []string) error {

	// TODO: if given just a list of unknown args, then assume they meant "sync"
	// TODO: if given --login, then force the login

	// loginFlag, _ := cmd.Flags().GetBool("login")

	cmd.Println("ROOT COMMAND")

	return nil
}
