//go:build !nodev

package devrunner

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func RunE(cmd *cobra.Command, args []string) error {

	switch args[0] {
	case "colors":
		return colorTest(cmd, args[1:])

	case "error":
		return errors.New("some error")
	}

	return nil
}

func colorTest(cmd *cobra.Command, _ []string) error {
	for i := 0; i < 300; i++ {
		out := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(i)).Render(fmt.Sprintf("This is color code #%d", i))
		cmd.Println(out)
	}
	return nil
}
