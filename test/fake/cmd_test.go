package fake

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tommyknows/packa/pkg/cmd"
)

func TestNoOpCmd(t *testing.T) {
	is := is.New(t)

	executedCommand := make(chan []string, 1)
	output := "testoutput"

	out, err := cmd.Execute([]string{"echo", "hello world"}, NoOp(executedCommand, output))
	is.Equal(output+"\n", out) // output of command should be expected plus newline
	is.NoErr(err)              // echo should not return an error

	is.Equal([]string{"echo", "hello world"}, <-executedCommand) // executed command should be echo hello world

	// set the global options
	cmd.AddGlobalOptions(NoOp(executedCommand, output))

	out, err = cmd.Execute([]string{"echo", "hello world"})
	is.Equal(output+"\n", out)
	is.NoErr(err)
	is.Equal([]string{"echo", "hello world"}, <-executedCommand)
	cmd.ResetGlobalOptions()
}

func TestNoOpErrorCmd(t *testing.T) {
	is := is.New(t)
	executedCommand := make(chan []string, 1)
	output := "testoutput"

	out, err := cmd.Execute([]string{"echo", "hello world"}, NoOpError(executedCommand, output))
	is.Equal(output+"\n", out)
	is.True(err != nil)
	is.Equal([]string{"echo", "hello world"}, <-executedCommand)
}
