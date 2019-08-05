package defaults

import (
	"os/user"
	"path"
)

const (
	packaHiddenDir = ".packa"
	configFileName = "packa.yml"
)

// WorkingDir returns the default working directory
func WorkingDir() string {
	usr, _ := user.Current()
	return path.Join(usr.HomeDir, packaHiddenDir)
}

// ConfigFileFullPath returns the full path to the
// configuration file
func ConfigFileFullPath() string {
	usr, _ := user.Current()
	return path.Join(usr.HomeDir, packaHiddenDir, configFileName)
}
