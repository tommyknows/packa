package gopackage

import (
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
)

const latest = "latest"

// Package contains info about a package that needs to be
// installed
type Package struct {
	// URL where to get the package from
	URL string `mapstructure:"url"`
	// InstallName to define the name of the binary
	InstallName string `mapstructure:"installName"`
	// Which version should be installed (semver, go modules!)
	Version string `mapstructure:"version"`
	// TODO
	// If the package should be auto-updated
	//AutoUpdate  bool
}

// Packages is a list of packages
type Packages []Package

// Add a package to package list if not already in it
func (pkgs *Packages) Add(pkg *Package) error {
	if pkgs.contain(pkg) {
		return errors.New("package already in list")
	}

	*pkgs = append(*pkgs, *pkg)
	return nil
}

func (pkgs *Packages) contain(pkg *Package) bool {
	for _, p := range *pkgs {
		if p.URL == pkg.URL && p.InstallName == pkg.InstallName {
			return true
		}
	}
	return false
}

// Install a given package
func (pkg *Package) Install() error {
	// change to ~/.packago directory to call go-commands
	// as there should not be a go.mod file
	cmd := exec.Command("go", "get", "-d", pkg.URL+"@"+pkg.Version)
	cmd.Dir = path.Join(os.Getenv("HOME"), ".packago")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Remove a given package
func (pkg *Package) Remove(autoremove bool) error {
	gopath := os.Getenv("GOPATH")
	err := os.Remove(path.Join(gopath, "bin", pkg.InstallName))
	if err != nil {
		return errors.Wrapf(err, "error removing binary %v", pkg.InstallName)
	}
	if !autoremove {
		return nil
	}
	err = os.Remove(path.Join(gopath, "src", pkg.URL))
	return errors.Wrapf(err, "error removing source code to binary %v", pkg.URL)
}
