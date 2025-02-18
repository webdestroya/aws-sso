package envrunner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/webdestroya/awssso/internal/runners/credentialsrunner"
)

func Run(cmd *cobra.Command, args []string) {
	_ = RunE(cmd, args)
}

func RunE(cmd *cobra.Command, args []string) error {

	// cmd.Printf("ARGS: %v\n", args)

	profile, args := args[0], args[1:]

	command, args := args[0], args[1:]

	if command == "--" {
		if len(args) < 2 {
			return errors.New("No command was provided")
		}
		command, args = args[0], args[1:]
	}

	// cmd.Printf("COMMAND: %s\n", command)
	// cmd.Printf("CARGS: %v\n", args)

	binPath, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	// cmd.Printf("BinPath: %s\n", binPath)

	credinfo, err := credentialsrunner.GetAWSCredentials(cmd.Context(), cmd.OutOrStdout(), profile)
	if err != nil {
		return err
	}
	// cmd.Printf("Got Creds\n")

	env := os.Environ()
	// env := make([]string, 0)

	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", credinfo.AccessKeyID))
	env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", credinfo.SecretAccessKey))
	env = append(env, fmt.Sprintf("AWS_SESSION_TOKEN=%s", credinfo.SessionToken))
	// env = append(env, fmt.Sprintf("AWS_REGION=%s", credinfo.Region))

	innerCmd := exec.CommandContext(cmd.Context(), binPath, args...)
	innerCmd.Stdin = os.Stdin
	innerCmd.Stdout = os.Stdout
	innerCmd.Stderr = os.Stderr
	innerCmd.Env = env
	innerCmd.Start()
	done := make(chan struct{})

	go func() {
		err := innerCmd.Wait()
		_ = err
		// status := innerCmd.ProcessState.Sys().(syscall.WaitStatus)
		// exitStatus := status.ExitStatus()
		// signaled := status.Signaled()
		// signal := status.Signal()
		close(done)
	}()
	// innerCmd.Process.Kill()
	<-done

	return nil

	// return syscall.Exec(binPath, args, env)
}
