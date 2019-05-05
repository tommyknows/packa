package fake

// CommandHandler executes commands
type CommandHandler struct {
	Output          string
	Err             error
	RemovedBinaries []string
	// returns the repo of the installed packages
	InstalledPackages []string
}

// NewCommandHandler creates a new handler for commands
func NewCommandHandler(output string, err error) *CommandHandler {
	// create a buffered channel so that we don't block
	return &CommandHandler{
		Output:            output,
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
	return h.Output, h.Err
}

// Remove a binary from gopath
func (h *CommandHandler) Remove(binary string) error {
	h.RemovedBinaries = append(h.RemovedBinaries, binary)
	return h.Err
}
