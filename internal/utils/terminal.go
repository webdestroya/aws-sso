package utils

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsTerminal determines if a file descriptor is an interactive terminal / TTY.
// first param is probably an io.Reader/io.Writer
func IsTerminal(v any) bool {
	f, ok := v.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
