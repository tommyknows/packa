package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"

	"k8s.io/klog"
)

const ()

// GoInstall calls go get to install a package
func GoInstall(repo string) (string, error) {
	// create buffer / writer for command output
	var b bytes.Buffer
	var stdoutMW, stderrMW io.Writer
	if klog.V(1) {
		stdoutMW = io.MultiWriter(&b, os.Stdout)
		stderrMW = io.MultiWriter(&b, os.Stderr)
	} else {
		stderrMW = io.Writer(&b)
		stderrMW = io.Writer(&b)
	}

	cmd := exec.Command("go", "get", repo)
	// TODO: get working dir from constants
	cmd.Dir = path.Join(os.Getenv("HOME"), ".packago")
	cmd.Stdout = stdoutMW
	cmd.Stderr = stderrMW
	err := cmd.Run()
	return b.String(), err
}
