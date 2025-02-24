//go:build !nodev

package devrunner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func RunE(cmd *cobra.Command, args []string) error {

	switch args[0] {
	case "colors":
		return colorTest(cmd, args[1:])

	case "spinner":
		return spinnerTest(cmd, args)

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

func spinnerTest(cmd *cobra.Command, _ []string) error {

	spinCtx, cancelFunc := context.WithCancel(cmd.Context())
	defer cancelFunc()

	spinChan := make(chan struct{})

	action := func(ctx context.Context) error {
		<-spinChan
		return nil
	}

	go func() {
		spin := spinner.New().Context(spinCtx).ActionWithErr(action).Output(cmd.ErrOrStderr()).Title("Waiting for stuff...")
		_ = spin.Run()
	}()

	time.Sleep(5 * time.Second)

	close(spinChan)
	fmt.Fprintln(cmd.ErrOrStderr(), "TIME UP")

	fmt.Fprintln(cmd.ErrOrStderr(), "DONE")
	return nil
}
