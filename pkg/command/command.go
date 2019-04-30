package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"k8s.io/klog"
)

// GoInstall calls go get to install a package and returns
// output and exit code
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
	cmd.Dir = "/Users/ramon/.packago/"
	cmd.Stdout = stdoutMW
	cmd.Stderr = stderrMW
	err := cmd.Run()
	return b.String(), err
}
