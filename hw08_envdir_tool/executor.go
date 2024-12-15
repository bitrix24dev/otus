package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) // #nosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	command.Env = os.Environ()
	for key, value := range env {
		if value.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=", key))
		} else {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value.Value))
		}
	}

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "failed to run command: %v\n", err)
		return 1
	}

	return 0
}
