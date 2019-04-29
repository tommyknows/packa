package packages

import (
	"fmt"
	"strings"

	"git.ramonruettimann.ml/ramon/packago/pkg/command"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

const (
	latest                  = "latest"
	master                  = "master"
	errWrapGoGet            = "Error executing go get"
	extractPackagePrefix    = "go: extracting "
	findingPackagePrefix    = "go: finding "
	errAlreadyDownloaded    = Error("Package has already been downloaded")
	errNoExtractPrefix      = Error("Output does not contain the 'go: extracting' prefix")
	errInvalidExtractOutput = Error("Go Get output does not contain extract dir")
	errNoGoDownloadOutput   = Error("Go download output is empty string")
	// ErrPackageAlreadyInstalled is returned if a package is already installed
	ErrPackageAlreadyInstalled   = Error("Package has already been installed")
	errNoUpdateNeeded            = Error("No update needed as version is pinned")
	errWrapInstallingAllPackages = "Error installing all packages"
	errWrapUpgradingAllPackages  = "Error upgrading all packages"
	errWrapUpgradePackageFailed  = "Error upgrading package"
)

// Error implements the error type in this package
type Error string

// Error implements error
func (e Error) Error() string { return string(e) }

// Package contains info about a package that needs to be
// installed
type Package struct {
	// URL where to get the package from
	URL string `yaml:"URL"`
	// Which version should be installed (semver, go modules!)
	Version string `yaml:"Version"`
	// internal: InstalledVersion
	InstalledVersion string `yaml:"InstalledVersion,omitempty"`
	// TODO
	// If the package should be auto-updated
	//AutoUpdate  bool
}

// Packages is a list of packages (that are installed!)
type Packages []*Package

// GetPackage returns the package identified by url
func (pkgs *Packages) GetPackage(url string) *Package {
	for _, p := range *pkgs {
		if p.URL == url {
			return p
		}
	}
	return nil
}

// UpgradeAll packages if needed
func (pkgs *Packages) UpgradeAll() error {
	for _, pkg := range *pkgs {
		err := pkg.UpgradeLatest()
		if err != nil && err != errNoUpdateNeeded {
			return errors.Wrapf(err, errWrapUpgradingAllPackages)
		}
	}
	return nil
}

// CreatePackage takes a URL and returns a package
func CreatePackage(url string) *Package {
	pkg := &Package{}
	lastIdx := strings.LastIndex(url, "@")
	// No version is given
	if lastIdx == -1 {
		pkg.URL = url
		pkg.Version = latest
		return pkg
	}

	pkg.URL = url[:lastIdx]
	pkg.Version = url[lastIdx+1:]
	return pkg
}

// Install a given package if not installed already,
// also add it to the list
func (pkgs *Packages) Install(pkg *Package) error {
	klog.V(3).Infof("Installed called on package %v", pkg)
	// package already in list, check if upgrade required
	if p := pkgs.GetPackage(pkg.URL); p != nil {
		klog.V(4).Infof("Comparing versions of %v and %v", p, pkg)
		if p.Version == pkg.Version {
			return ErrPackageAlreadyInstalled
		}
		klog.V(4).Infof("Changing package version from %v to %v", p.Version, pkg.Version)
		return p.UpgradeTo(pkg.Version)
	}

	err := pkg.install()
	if err != nil {
		return err
	}

	*pkgs = append(*pkgs, pkg)
	return nil
}

func (pkg *Package) getVersion(output string) string {
	if !strings.Contains(output, extractPackagePrefix) {
		return ""
	}

	var split []string
	for {
		split = strings.SplitN(output, "\n", 2)
		if strings.Contains(split[0], extractPackagePrefix+pkg.URL) {
			return split[0][strings.LastIndex(split[0], " ")+1:]
		}
		// last line
		if len(split) == 1 {
			return ""
		}
		output = split[1]
	}
}

// InstallAll packages
func (pkgs *Packages) InstallAll() error {
	for _, pkg := range *pkgs {
		// TODO: find out package name to check if binary already installed
		// and version of pkg.InstalledVersion is equal to pkg.Version
		err := pkg.install()
		if err != nil {
			return errors.Wrapf(err, errWrapInstallingAllPackages)
		}
	}
	return nil
}

func (pkgs *Packages) contain(pkg *Package) bool {
	for _, p := range *pkgs {
		if p.URL == pkg.URL {
			return true
		}
	}
	return false
}

// install a given package and set the installed
// version
func (pkg *Package) install() error {
	fmt.Printf("Installing Package %v@%v...\n", pkg.URL, pkg.Version)
	output, err := command.GoInstall(pkg.URL + "@" + pkg.Version)
	if err != nil {
		return err
	}

	version := pkg.getVersion(output)

	switch {
	case output == "":
		break
	case version != "":
		pkg.InstalledVersion = version
	case pkg.InstalledVersion == "":
		pkg.InstalledVersion = "~" + pkg.Version
	case !strings.HasPrefix(pkg.InstalledVersion, "~"):
		pkg.InstalledVersion = "~" + pkg.InstalledVersion
	default:
		// DO NOT TOUCH THE INSTALLED VERSION
	}

	fmt.Printf("Installed Package %v@%v\n", pkg.URL, pkg.Version)
	return nil
}

// UpgradeLatest upgrades a package if version is set
// to latest or master
func (pkg *Package) UpgradeLatest() error {
	if pkg.Version != latest && pkg.Version != master {
		return errNoUpdateNeeded
	}
	err := pkg.install()
	return err
}

// UpgradeTo a specific version and then set the package's version
// if installation was successful
func (pkg *Package) UpgradeTo(newVersion string) error {
	err := pkg.install()
	if err != nil {
		return errors.Wrapf(err, errWrapUpgradePackageFailed)
	}
	pkg.Version = newVersion
	return nil
}
