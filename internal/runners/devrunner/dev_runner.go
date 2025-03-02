//go:build !nodev

package devrunner

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/utils/cmdutils"
)

func RunE(cmd *cobra.Command, args []string) error {

	switch args[0] {
	case "colors":
		return colorTest(cmd, args[1:])

	case "error":
		return errors.New("some error")

	case "nonusage":
		return cmdutils.NewNonUsageError("some non usage error")
	}

	return nil
}

func colorTest(cmd *cobra.Command, _ []string) error {
	for i := range 300 {
		out := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(i)).Render(fmt.Sprintf("This is color code #%d", i))
		cmd.Println(out)
	}
	return nil
}
