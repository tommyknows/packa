package packages

import (
	"fmt"
	"strings"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

const (
	latest               = "latest"
	master               = "master"
	extractPackagePrefix = "go: extracting "
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

	fmt.Printf("LastIndex in URL %v: %v\n", url, lastIdx)
	pkg.Package.URL = url[:lastIdx]
	pkg.Package.Version = url[lastIdx+1:]
	return pkg
}

// Remove binary for a given package
func (pkg Package) Remove() error {
	fmt.Printf("Removing Package %v@%v...\n", pkg.URL, pkg.Version)
	lastIndex := strings.LastIndex(pkg.URL, "/")
	binaryName := pkg.URL[lastIndex+1:]
	err := pkg.cmdHandler.Remove(binaryName)
	if err != nil {
		return errors.Wrapf(err, "could not remove package")
	}
	fmt.Printf("Removed Package %s@%s\n", pkg.URL, pkg.Version)
	return nil
}

// Install a given package and set the installed
// version
func (pkg Package) Install() error {
	fmt.Printf("Installing Package %v@%v...\n", pkg.URL, pkg.Version)
	output, err := pkg.cmdHandler.Install(pkg.URL, pkg.Version)
	if err != nil {
		return errors.Wrapf(err, "could not install package %v", pkg.URL)
	}

	// print info at the end, anonymous function to have pkg.InstalledVersion set
	defer func() {
		fmt.Printf("Installed Package %s@%s\n", pkg.URL, pkg.InstalledVersion)
	}()

	version := pkg.getVersion(output)
	if version != "" {
		klog.V(1).Infof("Determined version from output: %v", version)
		pkg.InstalledVersion = version
		return nil
	}

	// we could not get version from go get output...
	// if the output was empty, it should've been the
	// version that was specified, i assume
	if output == "" {
		klog.V(1).Infof("No go get output on installation, setting version %v", pkg.Version)
		pkg.InstalledVersion = pkg.Version
		return nil
	}

	// we want to set something, so guess to specified version
	klog.V(1).Infof("Setting version as unsure to %v", pkg.Version)
	pkg.InstalledVersion = "~" + pkg.Version
	return nil
}

// UpgradeTo a specific version
func (pkg Package) UpgradeTo(newVersion string) error {
	pkg.Version = newVersion
	err := pkg.Install()
	return errors.Wrapf(err, errWrapUpgradePackageFailed)
}

func (pkg Package) getVersion(output string) string {
	if !strings.Contains(output, extractPackagePrefix) {
		return ""
	}

	var split []string
	for {
		split = strings.SplitN(output, "\n", 2)
		if strings.Contains(split[0], extractPackagePrefix+pkg.URL) {
			return split[0][strings.LastIndex(split[0], " ")+1:]
		}
		output = split[1]
	}
}
