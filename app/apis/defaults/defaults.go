package defaults

import (
	"os"
	"path"
)

const (
	packaHiddenDir = ".packa"
	configFileName = "packa.yml"
)

// WorkingDir returns the default morking directory
func WorkingDir() string {
	return path.Join(os.Getenv("HOME"), packaHiddenDir)
}

// ConfigDir returns the default directory
// where the config can be found
func ConfigDir() string {
	return path.Join(os.Getenv("HOME"), packaHiddenDir)
}

// ConfigFilename retuns the name of the configfile
func ConfigFilename() string {
	return configFileName
}

// ConfigFileFullPath returns the full path to the
// configuration file
func ConfigFileFullPath() string {
	return path.Join(ConfigDir(), ConfigFilename())
}

// BinaryDir returns the directory where the binaries
// will be deleted from on a remove-operation
func BinaryDir() string {
	return path.Join(os.Getenv("GOPATH"), "bin")
}
