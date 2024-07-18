package cmd

import "github.com/spf13/cobra"

func vPrintf(cmd *cobra.Command, strfmt string, args ...any) {
	if !verboseLogging {
		return
	}

	cmd.Printf(strfmt, args...)
}

func vPrintln(cmd *cobra.Command, args ...any) {
	if !verboseLogging {
		return
	}

	cmd.Println(args...)
}
