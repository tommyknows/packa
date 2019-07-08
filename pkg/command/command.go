package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"git.ramonruettimann.ml/ramon/packa/app/apis/defaults"
	"k8s.io/klog"
)

const (
	extractPackagePrefix = "go: extracting "
)

// Handler executes commands
type Handler struct {
	workingDir string
}

// NewHandler creates a new handler for commands
func NewHandler(opts ...func(*Handler) error) (*Handler, error) {
	h := &Handler{
		workingDir: defaults.WorkingDir(),
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

// InstallError is what is returned
type InstallError struct {
	// the actual error from the command
	err error
	// the commands output
	output string
}

// newInstallError returns an error of type install errror.
// if the given error was nil, a nil error will be returned
func newInstallError(output string, err error) error {
	if err == nil {
		return nil
	}
	return InstallError{
		err,
		output,
	}
}

// Cause returns the original error
func (e InstallError) Cause() error {
	return e.err
}

func (e InstallError) Error() string {
	// output normally contains newline already
	return fmt.Sprintf("%v%v", e.output, e.err)
}

// Install calls go get to install a package and returns
// output and exit code
func (h *Handler) Install(repo, version string) (extractedVersion string, err error) {
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

	cmd := exec.Command("go", "get", repo+"@"+version)
	cmd.Dir = h.workingDir
	cmd.Stdout = stdoutMW
	cmd.Stderr = stderrMW
	err = cmd.Run()
	output := b.String()
	if err != nil {
		return "", newInstallError(output, err)
	}

	extractedVersion = h.getVersion(repo, output)
	if extractedVersion != "" {
		klog.V(1).Infof("Determined version from output: %v", version)
		return extractedVersion, nil
	}

	// if the output was empty, it should've been the
	// version that was specified, i assume
	if output == "" {
		klog.V(1).Infof("No go get output on installation, setting version %v", version)
		return version, nil
	}
	klog.V(1).Infof("Setting version as unsure to %v", version)
	return "~" + version, nil
}

func (h *Handler) getVersion(repoURL, output string) string {
	if !strings.Contains(output, extractPackagePrefix+repoURL) {
		return ""
	}

	var split []string
	for {
		split = strings.SplitN(output, "\n", 2)
		if strings.Contains(split[0], extractPackagePrefix+repoURL) {
			return split[0][strings.LastIndex(split[0], " ")+1:]
		}
		output = split[1]
	}
}

// Remove a binary from gopath
func (h *Handler) Remove(binary string) error {
	return os.RemoveAll(path.Join(os.Getenv("GOPATH"), "bin", binary))
}
