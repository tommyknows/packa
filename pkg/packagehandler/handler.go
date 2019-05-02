package packages

import (
	"fmt"
	"strings"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"github.com/pkg/errors"
)

const (
	// packageNotInstalled is the index that is returned
	// if the package is not installed / in the list
	// of packages
	packageNotInstalled = -1
)

// PackageHandler is a handler for packages
type PackageHandler struct {
	packages   []Package
	cmdHandler CommandHandler
}

// CommandHandler contains commands that are used to install packages
type CommandHandler interface {
	// Install a go package
	Install(repo string) (string, error)
	// Remove a go package (or just its binary)
	Remove(binaryName string) error
}

// NewPackageHandler creates a new handler for packages
func NewPackageHandler(cmdH CommandHandler, opts ...func(*PackageHandler) error) (*PackageHandler, error) {
	h := &PackageHandler{
		cmdHandler: cmdH,
	}
	for _, option := range opts {
		err := option(h)
		if err != nil {
			return nil, err
		}
	}
	return h, nil
}

// SetCommandHandler to use for installing and removing packages
func SetCommandHandler(cmdH CommandHandler) func(*PackageHandler) error {
	return func(pkgH *PackageHandler) error {
		pkgH.cmdHandler = cmdH
		return nil
	}
}

// Handle the packages given
func Handle(pkgs []*config.Package) func(*PackageHandler) error {
	return func(pkgH *PackageHandler) error {
		for _, p := range pkgs {
			pkgH.packages = append(pkgH.packages, Package{p, pkgH.cmdHandler})
		}
		return nil
	}
}

// ExportPackages as config.Package type
func (pkgH *PackageHandler) ExportPackages() []*config.Package {
	var pkgs []*config.Package
	for _, p := range pkgH.packages {
		pkgs = append(pkgs, p.Package)
	}

	return pkgs
}

// GetPackages returns the package identified by url or a new package if
// it did not exist yet
func (pkgH *PackageHandler) GetPackages(urls ...string) []Package {
	var pkgs []Package
	for _, url := range urls {
		pkgs = append(pkgs, pkgH.getPackage(url))
	}

	return pkgs
}

func (pkgH *PackageHandler) getPackage(url string) Package {
	var version string
	lastIdx := strings.LastIndex(url, "@")
	// No version is given
	if lastIdx == -1 {
		version = latest
	} else {
		version = url[lastIdx+1:]
		url = url[:lastIdx]
	}

	for _, p := range pkgH.packages {
		if p.URL == url {
			return p
		}
	}
	return Package{
		&config.Package{
			URL:     url,
			Version: version,
		},
		pkgH.cmdHandler,
	}
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
		if err != nil {
			return errors.Wrapf(err, errWrapUpgradingAllPackages)
		}
	}
	return nil
}

// InstallError "collects" all errors while installing package(s)
type InstallError map[Package]error

// Add an error if non-nil
func (ie InstallError) add(pkg Package, err error) {
	if err != nil {
		ie[pkg] = err
	}
}

func (ie InstallError) Error() string {
	s := "Encountered error(s) while installing packages:\n"
	for pkg, err := range ie {
		s += fmt.Sprintf("Package %s@%s: %s\n", pkg.URL, pkg.Version, err.Error())
	}

	return s
}

// Install the given packages and add them to the
// list of packages. Returns an InstallError that contains
// a map of the failed packages withe the error message.
// if an error occurs, the
func (pkgH *PackageHandler) Install(pkgs ...Package) error {
	// contain is just a variable to optimise the code
	// as there is no need to check if the package is in
	// the list of packages if we just took it from there
	contain := false
	// if no package is specified, go through them all
	if len(pkgs) == 0 {
		pkgs = pkgH.packages
		contain = true
	}

	collectionErr := make(InstallError)
	for _, pkg := range pkgs {
		err := pkg.Install()
		if err != nil {
			collectionErr.add(pkg, err)
			continue
		}
		if !contain && pkgH.index(pkg) == packageNotInstalled {
			pkgH.packages = append(pkgH.packages, pkg)
		}
	}
	return nil
}

// Remove binaries and from the list
func (pkgH *PackageHandler) Remove(pkgs ...Package) error {
	collectionErr := make(InstallError)
	for _, pkg := range pkgs {
		idx := pkgH.index(pkg)
		if pkgH.index(pkg) == packageNotInstalled {
			//rn errors.Errorf("package %v not installed", pkg.URL)
			collectionErr.add(pkg, errors.Errorf("package %v not installed", pkg.URL))
			continue
		}

		err := pkg.Remove()
		if err != nil {
			collectionErr.add(pkg, errors.Wrapf(err, "error removing binary, not removing package %v from state file", pkg.URL))
			continue
		}

		pkgH.packages = append(pkgH.packages[:idx], pkgH.packages[idx+1:]...)
	}
	return collectionErr
}

// index returns the index of the package or -1 (packageNotInstalled)
func (pkgH *PackageHandler) index(pkg Package) int {
	for i, p := range pkgH.packages {
		if p.URL == pkg.URL {
			return i
		}
	}
	return packageNotInstalled
}
