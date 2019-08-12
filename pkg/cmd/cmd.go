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

type Cmd *exec.Cmd

type Option func(Cmd) error

// globalOpts define options that will be added
// to every cmd.Exec call
var globalOpts []Option

func AddGlobalOptions(opts ...Option) {
	globalOpts = append(globalOpts, opts...)
}

func ResetGlobalOptions() {
	globalOpts = nil
}

// Execute is a kind-of simplified version of exec.Cmd with functional
// options.
func Execute(args []string, opts ...Option) (output string, err error) {
	var c Cmd
	if len(args) == 0 {
		return "", errors.New("no arguments / command given")
	}
	c = exec.Command(args[0], args[1:]...)

	var b bytes.Buffer
	c.Stdout, c.Stderr = io.Writer(&b), io.Writer(&b)

	for _, opt := range globalOpts {
		if err := opt(c); err != nil {
			return "", errors.Wrapf(err, "could not set global option")
		}
	}

	for _, opt := range opts {
		if err = opt(c); err != nil {
			return "", errors.Wrapf(err, "could not set option")
		}
	}

	err = (*c).Run()
	return b.String(), err
}

// WorkingDir sets the workingDirectory of the command, so
// in which directory the command will be executed
func WorkingDir(wd string) Option {
	return func(c Cmd) (err error) {
		c.Dir, err = expand(wd)
		return err
	}
}

// DirectPrint prints the output of the command to stdout / stderr
// if b is true.
func DirectPrint(b bool) Option {
	return func(c Cmd) error {
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
