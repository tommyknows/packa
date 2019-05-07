package fake

import "fmt"

// CommandHandler executes commands
type CommandHandler struct {
	Output          string
	Version         string
	Err             error
	RemovedBinaries []string
	// returns the repo of the installed packages
	InstalledPackages []string
}

// NewCommandHandler creates a new handler for commands
func NewCommandHandler(output, version string, err error) *CommandHandler {
	// create a buffered channel so that we don't block
	return &CommandHandler{
		Output:            output,
		Version:           version,
		Err:               err,
		RemovedBinaries:   []string{},
		InstalledPackages: []string{},
	}
}

// Install calls go get to install a package and returns
// output and exit code
func (h *CommandHandler) Install(repo, version string) (string, error) {
	// create buffer / writer for command output
	h.InstalledPackages = append(h.InstalledPackages, repo)
	if h.Err != nil {
		return h.Version, fmt.Errorf("%v%v", h.Output, h.Err)
	}
	return h.Version, nil
}

// Remove a binary from gopath
func (h *CommandHandler) Remove(binary string) error {
	h.RemovedBinaries = append(h.RemovedBinaries, binary)
	return h.Err
}
