package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCommand executes a command in the specified directory.
func RunCommand(cwd string, cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Dir = cwd
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Println("Running command: ", command.String())
	if err := command.Run(); err != nil {
		return "", err
	}
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil

}
