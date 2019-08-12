package fake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommyknows/packa/pkg/cmd"
)

func TestNoOpCmd(t *testing.T) {
	executedCommand := make(chan []string, 1)
	output := "testoutput"

	out, err := cmd.Execute([]string{"echo", "hello world"}, NoOp(executedCommand, output))
	assert.Equal(t, output+"\n", out)
	assert.NoError(t, err)
	assert.Equal(t, []string{"echo", "hello world"}, <-executedCommand)

	// set the global options
	cmd.AddGlobalOptions(NoOp(executedCommand, output))
	//cmd.GlobalOpts = []cmd.Option{NoOp(executedCommand, output)}

	out, err = cmd.Execute([]string{"echo", "hello world"})
	assert.Equal(t, output+"\n", out)
	assert.NoError(t, err)
	assert.Equal(t, []string{"echo", "hello world"}, <-executedCommand)
	cmd.ResetGlobalOptions()
}

func TestNoOpErrorCmd(t *testing.T) {
	executedCommand := make(chan []string, 1)
	output := "testoutput"

	out, err := cmd.Execute([]string{"echo", "hello world"}, NoOpError(executedCommand, output))
	assert.Equal(t, output+"\n", out)
	assert.Error(t, err)
	assert.Equal(t, []string{"echo", "hello world"}, <-executedCommand)
}
