package packages

import (
	"strings"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	"github.com/pkg/errors"
)

const (
	latest = "latest"
	master = "master"
	// ErrPackageAlreadyInstalled is returned if a package is already installed
	ErrPackageAlreadyInstalled   = Error("Package has already been installed")
	errWrapInstallingAllPackages = "Error installing all packages"
	errWrapUpgradingAllPackages  = "Error upgrading all packages"
	errWrapUpgradePackageFailed  = "Error upgrading package"
)

// Error implements the error type in this package
type Error string

// Error implements error
func (e Error) Error() string { return string(e) }

// Package is a wrapper for the config package
type Package struct {
	*config.Package
	cmdHandler CommandHandler
}

// NewPackage takes a URL and returns a package
func NewPackage(url string, cmdH CommandHandler) Package {
	pkg := Package{
		Package:    &config.Package{},
		cmdHandler: cmdH,
	}

	lastIdx := strings.LastIndex(url, "@")
	// No version is given
	if lastIdx == -1 {
		pkg.Package.URL = url
		pkg.Package.Version = latest
		return pkg
	}

	pkg.Package.URL = url[:lastIdx]
	pkg.Package.Version = url[lastIdx+1:]
	return pkg
}

// Remove binary for a given package
func (pkg Package) Remove() error {
	output.Info("📦 Removing %v@%v...\n", pkg.URL, pkg.Version)
	lastIndex := strings.LastIndex(pkg.URL, "/")
	binaryName := pkg.URL[lastIndex+1:]
	err := pkg.cmdHandler.Remove(binaryName)
	if err != nil {
		return errors.Wrapf(err, "could not remove package")
	}
	output.Success("📦 Removed %s@%s", pkg.URL, pkg.InstalledVersion)
	return nil
}

// Install a given package and set the installed
// version
func (pkg Package) Install() error {
	output.Info("📦 Installing %v@%v...\n", pkg.URL, pkg.Version)
	version, err := pkg.cmdHandler.Install(pkg.URL, pkg.Version)
	if err != nil {
		return errors.Wrapf(err, "could not install package %v", pkg.URL)
	}

	pkg.InstalledVersion = version
	output.Success("📦 Installed %s@%s", pkg.URL, pkg.InstalledVersion)
	return nil
}

// UpgradeTo a specific version
func (pkg Package) UpgradeTo(newVersion string) error {
	pkg.Version = newVersion
	err := pkg.Install()
	return errors.Wrapf(err, errWrapUpgradePackageFailed)
}
