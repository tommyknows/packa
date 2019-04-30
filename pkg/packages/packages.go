package packages

import (
	"fmt"
	"os"
	"path"
	"strings"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
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
	errNoUpgradeNeeded           = Error("No update needed as version is pinned")
	errWrapInstallingAllPackages = "Error installing all packages"
	errWrapUpgradingAllPackages  = "Error upgrading all packages"
	errWrapUpgradePackageFailed  = "Error upgrading package"
)

// Error implements the error type in this package
type Error string

// Error implements error
func (e Error) Error() string { return string(e) }

// PackageHandler is a handler for packages
type PackageHandler struct {
	packages   []Package
	workingDir string
}

// Package is a wrapper for the config package
type Package struct {
	*config.Package
}

// NewPackageHandler creates a new handler for packages
func NewPackageHandler(opts ...func(*PackageHandler) error) (*PackageHandler, error) {
	h := &PackageHandler{}
	for _, option := range opts {
		err := option(h)
		if err != nil {
			return nil, err
		}
	}
	return h, nil
}

// Handle the packages given
func Handle(pkgs []*config.Package) func(*PackageHandler) error {
	return func(p *PackageHandler) error {
		return p.handle(pkgs)
	}
}

func (pkgH *PackageHandler) handle(pkgs []*config.Package) error {
	for _, p := range pkgs {
		pkgH.packages = append(pkgH.packages, Package{p})
	}

	return nil
}

// WorkingDir sets the working directory for the package handler
func WorkingDir(dir string) func(*PackageHandler) error {
	return func(p *PackageHandler) error {
		return p.setWorkingDir(dir)
	}
}

func (pkgH *PackageHandler) setWorkingDir(dir string) error {
	pkgH.workingDir = dir
	return nil
}

// ExportPackages as config.Package type
func (pkgH *PackageHandler) ExportPackages() []*config.Package {
	var pkgs []*config.Package
	for _, p := range pkgH.packages {
		pkgs = append(pkgs, p.Package)
	}

	return pkgs
}

// GetPackage returns the package identified by url
func (pkgH *PackageHandler) GetPackage(url string) Package {
	for _, p := range pkgH.packages {
		if p.URL == url {
			return p
		}
	}
	return Package{}
}

// UpgradeAll packages if needed
func (pkgH *PackageHandler) UpgradeAll() error {
	for _, pkg := range pkgH.packages {
		// no automatic upgrade if version is pinned to specific semver tag
		if pkg.Version != latest && pkg.Version != master {
			fmt.Printf("Not upgrading %v as pinned to %v\n", pkg.URL, pkg.Version)
			continue
		}
		err := pkg.Install()
		if err != nil && err != errNoUpgradeNeeded {
			return errors.Wrapf(err, errWrapUpgradingAllPackages)
		}
	}
	return nil
}

// CreatePackage takes a URL and returns a package
func CreatePackage(url string) Package {
	pkg := Package{&config.Package{}}
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
func (pkgH *PackageHandler) Install(pkg Package) error {
	klog.V(3).Infof("Installed called on package %v", pkg)
	// package already in list, check if upgrade required
	if p := pkgH.GetPackage(pkg.URL); p.Package != nil {
		klog.V(4).Infof("Comparing versions of %v and %v", p, pkg)
		if p.Version == pkg.Version {
			return ErrPackageAlreadyInstalled
		}
		klog.V(4).Infof("Changing package version from %v to %v", p.Version, pkg.Version)
		return p.UpgradeTo(pkg.Version)
	}

	err := pkg.Install()
	if err != nil {
		return err
	}

	pkgH.packages = append(pkgH.packages, pkg)
	return nil
}

// Remove a binary and remove it from the list
func (pkgH *PackageHandler) Remove(pkg Package) error {
	err := pkg.Remove()
	if err != nil {
		return errors.Wrapf(err, "error removing binary")
	}

	var i int
	var p Package
	for i, p = range pkgH.packages {
		if p.URL == pkg.URL {
			break
		}
	}
	pkgH.packages = append(pkgH.packages[:i], pkgH.packages[i+1:]...)
	return nil
}

// Remove binary for a given package
func (pkg Package) Remove() error {
	lastIndex := strings.LastIndex(pkg.URL, "/")
	binaryName := pkg.URL[lastIndex+1:]
	return os.RemoveAll(path.Join(os.Getenv("GOPATH"), "bin", binaryName))
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
		// last line
		if len(split) == 1 {
			return ""
		}
		output = split[1]
	}
}

// InstallAll packages
func (pkgH *PackageHandler) InstallAll() error {
	for _, pkg := range pkgH.packages {
		// TODO: find out package name to check if binary already installed
		// and version of pkg.InstalledVersion is equal to pkg.Version
		err := pkg.Install()
		if err != nil {
			return errors.Wrapf(err, errWrapInstallingAllPackages)
		}
	}
	return nil
}

func (pkgH *PackageHandler) contain(pkg Package) bool {
	for _, p := range pkgH.packages {
		if p.URL == pkg.URL {
			return true
		}
	}
	return false
}

// Install a given package and set the installed
// version
func (pkg Package) Install() error {
	fmt.Printf("Installing Package %v@%v...\n", pkg.URL, pkg.Version)
	output, err := command.GoInstall(pkg.URL + "@" + pkg.Version)
	if err != nil {
		// output already contains newline
		return fmt.Errorf("%v%v", output, err)
	}

	version := pkg.getVersion(output)

	// TODO: what?
	switch {
	case output == "":
		pkg.InstalledVersion = pkg.Version
	case version != "":
		pkg.InstalledVersion = version
	case pkg.InstalledVersion == "":
		pkg.InstalledVersion = "~" + pkg.Version
	case strings.HasPrefix(pkg.InstalledVersion, "~"):
		break
	default:
		pkg.InstalledVersion = "~" + pkg.Version
	}

	fmt.Printf("Installed Package %v@%v\n", pkg.URL, pkg.InstalledVersion)
	return nil
}

// UpgradeTo a specific version and then set the package's version
// if installation was successful
func (pkg Package) UpgradeTo(newVersion string) error {
	pkg.Version = newVersion
	err := pkg.Install()
	if err != nil {
		return errors.Wrapf(err, errWrapUpgradePackageFailed)
	}
	return nil
}
