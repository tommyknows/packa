package fake

import (
	"fmt"

	"github.com/tommyknows/packa/pkg/cmd"
)

// NoOp makes cmd.Exec not execute the actual command, but rather just print the
// given output (adding a newline). Writes the command that would have been executed
// into the channel
func NoOp(cmds chan []string, output string) func(cmd.Cmd) error {
	return func(command cmd.Cmd) error {
		cmds <- command.Args
		command.Args = []string{"echo", output}
		command.Path = "/bin/echo"
		return nil
	}
}

// NoOpError acts the same as NooOp, but will make the command exit with a non-zero
// exit code.
func NoOpError(cmds chan []string, output string) func(cmd.Cmd) error {
	return func(command cmd.Cmd) error {
		cmds <- command.Args
		command.Args = []string{"sh", "-c", fmt.Sprintf("echo \"%v\" && false", output)}
		command.Path = "/bin/sh"
		return nil
	}
}
