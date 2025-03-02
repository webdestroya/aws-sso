//go:build !windows
// +build !windows

package envrunner

import (
	"os/exec"
	"syscall"
)

func ensureChildProcessesAreKilled(opts *envOptions, cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.Cancel = func() error {
		return syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}
	cmd.WaitDelay = opts.WaitDelay
}
