package goget

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/pkg/collection"
	"github.com/tommyknows/packa/pkg/defaults"
	"github.com/tommyknows/packa/pkg/output"
	"k8s.io/klog"
)

var (
	// taken from https://github.com/semver/semver/issues/232 with leading 'v'
	semVerRegex = regexp.MustCompile(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`)
	// matches if a (go get) URL contains a major version in the end
	majorVersionRegex = regexp.MustCompile(`/v([0-9])`)
)

type Handler struct {
	Config   configuration
	Packages []Package
}

type configuration struct {
	// Define in which directory the go get command should be executed
	WorkingDir string `json:"workingDir,omitempty"`
	// If go get should be called with the "-u" flag to update dependencies
	UpdateDependencies bool `json:"updateDependencies,omitempty"`
	// Print the output of the go get command
	PrintCommandOutput bool `json:"printCommandOutput,omitempty"`
}

type Package struct {
	URL     string `json:"url"`
	Version string `json:"version,omitempty"`
}

func (p Package) String() string {
	s := p.URL
	if p.Version != "" {
		s += "@" + p.Version
	}
	return s
}

// Init initialises the handler. If the packagelist should be nil, it adds itself
// to that list
func (goH *Handler) Init(config *json.RawMessage, packages *json.RawMessage) error {
	if config != nil {
		err := json.Unmarshal([]byte(*config), &goH.Config)
		if err != nil {
			return errors.Wrapf(err, "could not parse config %s", config)
		}
		klog.V(4).Infof("GoGet: unmarshaled config %#v", goH.Config)
	}

	if packages == nil {
		klog.V(4).Infof("GoGet: no packages defined, adding default package")
		goH.Packages = []Package{
			{
				URL:     "github.com/tommyknows/packa",
				Version: "latest",
			},
		}
		return nil
	}

	err := json.Unmarshal([]byte(*packages), &goH.Packages)
	if err != nil {
		return errors.Wrapf(err, "could not parse packages %s", packages)
	}
	klog.V(4).Infof("GoGet: Added package list %v", goH.Packages)
	return nil
}

// New returns a handler with the default settings. They will be overwritten
// (if set) on Init()
func New() *Handler {
	return &Handler{
		// the default handler config
		Config: configuration{
			WorkingDir:         defaults.WorkingDir(),
			UpdateDependencies: false,
			PrintCommandOutput: false,
		},
	}
}

// Install the packages and add them to the index. If an error occurs while installing
// a package, the other packages will be handled / installed nonetheless
func (goH *Handler) Install(pkgs ...string) (packageList *json.RawMessage, err error) {
	return goH.do(goH.install, goH.addToIndex, pkgs...)
}

// Remove a package from the system by parsing the given name and finding out the
// binary name. then remove the package from the index
func (goH *Handler) Remove(pkgs ...string) (packageList *json.RawMessage, err error) {
	return goH.do(goH.remove, goH.removeFromIndex, pkgs...)
}

// Upgrade a package, if it is in the index. Returns an error if a package
// should not exist in the index, but still processes all other packages
func (goH *Handler) Upgrade(pkgs ...string) (packageList *json.RawMessage, err error) {
	return goH.do(goH.upgrade, goH.upgradeIndex, pkgs...)
}

// do the package action and indexAction for a list of packages, handling
// errors and marshaling the index in the end
func (goH *Handler) do(packageAction func(Package) error, indexAction func(Package), pkgs ...string) (*json.RawMessage, error) {
	var pError collection.Error
	packages, err := goH.getPackages(pkgs...)
	if err != nil {
		pe, ok := err.(*collection.Error)
		if !ok {
			// we always expect a collection.Error from getPackages
			return nil, errors.Wrapf(err, "unexpected error occured when getting package list")
		}
		pError.Merge(*pe)
	}

	// should never occur, but as a safety check
	if packageAction == nil {
		return nil, errors.New("no package action defined")
	}

	for _, p := range packages {
		// execute the packageAction and then the index action, if applicable
		switch err := packageAction(p); {
		case indexAction == nil:
			klog.V(5).Infof("GoGet: Not executing index action because none is defined")
		case err != nil:
			klog.V(4).Infof("GoGet: Error while executing package action for %v, adding error to collection", p.String())
			pError.Add(p.String(), err)
			klog.V(5).Infof("GoGet: Not executing index action because of error on package action")
		default:
			indexAction(p)
		}
	}

	klog.V(6).Infof("GoGet: Marshaling packages")
	raw, err := json.Marshal(goH.Packages)
	if err != nil {
		return nil, err
	}
	msg := json.RawMessage(raw)
	return &msg, pError.IfNotEmpty()
}

// has package in index with version match
func (goH *Handler) has(p Package) bool {
	for _, pkg := range goH.Packages {
		if pkg == p {
			return true
		}
	}
	return false
}

// has package in index without version match
func (goH *Handler) hasURL(p Package) bool {
	for _, pkg := range goH.Packages {
		if pkg.URL == p.URL {
			return true
		}
	}
	return false
}

// getPackages returns either a parsed list of the given packages
// or all the packages defined in the index, if no argument is given.
// no matter what is returned, the slice is save to use and modify.
// it can return an error for some packages, but always tries to parse
// all packages. All successfully parsed packages will be added to the
// slice, all others will be described in the error
func (goH *Handler) getPackages(pkgs ...string) ([]Package, error) {
	var e collection.Error
	if len(pkgs) == 0 {
		klog.V(5).Infof("GoGet: No packages defined, using all packages from go handler")
		// make it safe to modify
		packages := make([]Package, len(goH.Packages))
		copy(packages, goH.Packages)
		return packages, nil
	}

	var packages []Package
	for _, pkg := range pkgs {
		p, err := parse(pkg)
		if err != nil {
			klog.V(6).Infof("GoGet: Error when parsing %v: %v", pkg, err)
			e.Add(pkg, err)
			continue
		}
		klog.V(6).Infof("GoGet: Successfully parsed package %v (%s)", pkg, p)
		packages = append(packages, p)
	}
	return packages, e.IfNotEmpty()
}

// add a package to the index, overwrite already existing package's version
func (goH *Handler) addToIndex(p Package) {
	klog.V(6).Infof("GoGet: Adding package %s to index", p)
	for idx, pkg := range goH.Packages {
		if pkg.URL == p.URL {
			goH.Packages[idx] = p
			klog.V(6).Infof("GoGet: Package %s already in index, overwriting old version %v", p, pkg.Version)
			return
		}
	}
	goH.Packages = append(goH.Packages, p)
}

// remove a package from the index if the URL and version are equal
func (goH *Handler) removeFromIndex(p Package) {
	klog.V(6).Infof("GoGet: Removing package %s from index", p)
	for idx, pkg := range goH.Packages {
		if pkg == p {
			goH.Packages = append(goH.Packages[:idx], goH.Packages[idx+1:]...)
			klog.V(6).Infof("GoGet: Package %s removed from index", p)
			return
		}
	}
	klog.V(6).Infof("GoGet: Package %s has not been defined in index", p)
}

// upgradeIndex adds packages to the index only
// if they existed in the index already, thus updating
// the version of the package in the index
func (goH *Handler) upgradeIndex(p Package) {
	if goH.hasURL(p) {
		goH.addToIndex(p)
	}
}

// goGet the given package
func (goH *Handler) goGet(pkg Package) error {
	c := []string{"go", "get", pkg.String()}
	if goH.Config.UpdateDependencies {
		// insert "-u"
		c = append(c[:2], append([]string{"-u"}, c[2:]...)...)
	}

	if goH.Config.PrintCommandOutput {
		output.Info("%v", c)
	}

	out, err := cmd.Execute(
		c,
		cmd.WorkingDir(goH.Config.WorkingDir),
		cmd.DirectPrint(bool(klog.V(5)) || goH.Config.PrintCommandOutput),
	)

	// don't print the output twice if we have verbosity
	if err != nil && !bool(klog.V(5)) || goH.Config.PrintCommandOutput {
		output.Warn(out)
	}
	return err
}

// install the package. Does not add it to the package list
func (goH *Handler) install(pkg Package) error {
	output.Info("📦 GoGet\tInstalling Package %s", pkg)
	err := goH.goGet(pkg)
	if err != nil {
		return err
	}
	output.Success("📦 GoGet\tInstalled Package %s", pkg)
	return nil
}

// remove a package from the system. As this is kind-of guesswork (parsing
// the name of the binary), it will ask the user for confirmation
func (goH *Handler) remove(pkg Package) error {
	output.Info("📦 GoGet\tRemoving Package %s", pkg)
	binName := extractBinaryName(pkg.URL)
	e := exec.Command("go", "env", "GOPATH")
	out, err := e.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "could not determine GOPATH")
	}

	// remove newline and add /bin/
	dirName := string(out[:len(out)-1]) + "/bin/"
	confirmed := output.WithConfirmation("removing binary %s (%s)", binName, dirName+binName)
	if !confirmed {
		klog.V(5).Infof("GoGet: Binary removal not confirmed by user, aborting")
		return nil
	}

	err = os.Remove(dirName + binName)
	if err != nil {
		return errors.Wrapf(err, "could not delete %v", dirName+binName)
	}
	output.Success("📦 GoGet\tRemoved Package %s", pkg)
	return nil
}

// upgrade only "installs" a package if it is defined in the index
func (goH *Handler) upgrade(pkg Package) error {
	output.Info("📦 GoGet\tUpgrading Package %s", pkg)
	if !goH.hasURL(pkg) {
		return errors.Errorf("package %v not in index", pkg.String())
	}

	// we don't update if the version specified is a valid
	// semver version, meaning it is "pinned".
	if matchSemVer(pkg.URL) {
		output.Info("Not upgrading %v as version is pinned to %v", pkg.URL, pkg.Version)
		return nil
	}
	if err := goH.goGet(pkg); err != nil {
		return err
	}
	output.Success("📦 GoGet\tUpgraded Package %s", pkg)
	return nil
}

// returns true if supplied version is a valid semver
func matchSemVer(version string) bool {
	return semVerRegex.MatchString(version)
}

// extracts the binary name from a URL.
// as the URL _can_ contain the major version, strip
// that away if it should exist
func extractBinaryName(url string) string {
	url = strings.TrimRight(url, "/")

	if !majorVersionRegex.MatchString(url) {
		s := strings.Split(url, "/")
		return s[len(s)-1]
	}

	b := majorVersionRegex.Split(url, 2)
	s := strings.Split(b[0], "/")
	return s[len(s)-1]
}

// parse the given package string into a URL and version.
// if the version should not be specified in the string,
// it returns go get's defaults "latest"
func parse(pkg string) (Package, error) {
	if !strings.Contains(pkg, "@") {
		return Package{pkg, ""}, nil
	}
	sp := strings.Split(pkg, "@")
	if len(sp) > 2 {
		return Package{}, fmt.Errorf("invalid url %v", pkg)
	}
	return Package{sp[0], sp[1]}, nil
}
