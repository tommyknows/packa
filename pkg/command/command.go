package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"

	"k8s.io/klog"
)

// Handler executes commands
type Handler struct {
	workingDir string
}

// NewHandler creates a new handler for commands
func NewHandler(opts ...func(*Handler) error) (*Handler, error) {
	h := &Handler{
		// TODO: get sensible default value from config / constants
		workingDir: "/Users/ramon/.packago",
	}

	for _, option := range opts {
		err := option(h)
		if err != nil {
			return nil, err
		}
	}
	return h, nil
}

// WorkingDir sets the working directory of the command handler
func WorkingDir(workDir string) func(*Handler) error {
	return func(h *Handler) error {
		h.workingDir = workDir
		return nil
	}
}

// Install calls go get to install a package and returns
// output and exit code
func (h *Handler) Install(repo string) (string, error) {
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
	cmd.Dir = h.workingDir
	cmd.Stdout = stdoutMW
	cmd.Stderr = stderrMW
	err := cmd.Run()
	return b.String(), err
}

// Remove a binary from gopath
func (h *Handler) Remove(binary string) error {
	return os.RemoveAll(path.Join(os.Getenv("GOPATH"), "bin", binary))
}
