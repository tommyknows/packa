package constants

import (
	"os"
	"path"
)

const (
	// ConfigFileLocalDir is the user-local directory to search for
	// a config file
	ConfigFileLocalDir = "$HOME/.packago"
	// ConfigFileName is the default name for the config file
	ConfigFileName = "packago.yml"
)

// GetDefaultWorkingDir gets the default working directory
func GetDefaultWorkingDir() string {
	return path.Join(os.Getenv("HOME"), ".packago")
}
