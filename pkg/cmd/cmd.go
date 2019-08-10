// Package cmd executes commands on the shell
package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

type cmd *exec.Cmd

type option func(cmd) error

// Execute is a kind-of simplified version of exec.Cmd with functional
// options.
func Execute(args []string, opts ...option) (output string, err error) {
	var c cmd
	if len(args) == 0 {
		return "", errors.New("no arguments / command given")
	}
	c = exec.Command(args[0], args[1:]...)

	var b bytes.Buffer
	c.Stdout, c.Stderr = io.Writer(&b), io.Writer(&b)

	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return "", errors.Wrapf(err, "could not set option")
		}
	}

	err = (*c).Run()
	return b.String(), err
}

// WorkingDir sets the workingDirectory of the command, so
// in which directory the command will be executed
func WorkingDir(wd string) option {
	return func(c cmd) (err error) {
		c.Dir, err = expand(wd)
		return err
	}
}

// DirectPrint prints the output of the command to stdout / stderr
// if b is true.
func DirectPrint(b bool) option {
	return func(c cmd) error {
		if b {
			c.Stdout = io.MultiWriter(c.Stdout, os.Stdout)
			c.Stderr = io.MultiWriter(c.Stderr, os.Stderr)
		}
		return nil
	}
}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
