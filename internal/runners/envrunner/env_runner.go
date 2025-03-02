package envrunner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webdestroya/aws-sso/internal/helpers/getcreds"
)

// some code for the runner taken from dagger

func (opts *envOptions) runE(cmd *cobra.Command, args []string) error {

	errOut := cmd.ErrOrStderr()

	profile, command, args := args[0], args[1], args[2:]

	if command == "--" {
		if len(args) == 0 {
			return errors.New("No command was provided")
		}
		command, args = args[0], args[1:]
	}

	binPath, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	if !opts.Quiet {
		fmt.Fprintln(errOut, "Profile:", profile)
		fmt.Fprintln(errOut, "Command: (unquoted)")
		fmt.Fprintln(errOut, binPath, strings.Join(args, " "))
		fmt.Fprintln(errOut)
	}

	credinfo, err := getcreds.GetAWSCredentials(cmd.Context(), cmd.ErrOrStderr(), profile)
	if err != nil {
		return err
	}

	env := os.Environ()

	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", credinfo.AccessKeyID))
	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", credinfo.AccessKeyID))
	env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", credinfo.SecretAccessKey))
	env = append(env, fmt.Sprintf("AWS_SESSION_TOKEN=%s", credinfo.SessionToken))
	// AWS_CREDENTIAL_EXPIRATION
	// env = append(env, fmt.Sprintf("AWS_REGION=%s", credinfo.Region))

	subCmd := exec.CommandContext(cmd.Context(), binPath, args...)
	subCmd.Stdin = os.Stdin
	subCmd.Stdout = os.Stdout
	subCmd.Stderr = os.Stderr
	subCmd.Env = env

	// NB: go run lets its child process roam free when you interrupt it, so
	// make sure they all get signalled. (you don't normally notice this in a
	// shell because Ctrl+C sends to the process group.)
	ensureChildProcessesAreKilled(opts, subCmd)

	return subCmd.Run()

	// subCmd.Start()

	// done := make(chan struct{})

	// go func() {
	// 	err := subCmd.Wait()
	// 	_ = err
	// 	// status := subCmd.ProcessState.Sys().(syscall.WaitStatus)
	// 	// exitStatus := status.ExitStatus()
	// 	// signaled := status.Signaled()
	// 	// signal := status.Signal()
	// 	close(done)
	// }()
	// // subCmd.Process.Kill()
	// <-done

	// return nil

	// return syscall.Exec(binPath, args, env)
}
