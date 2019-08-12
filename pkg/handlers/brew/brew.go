package brew

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/pkg/collection"
	"github.com/tommyknows/packa/pkg/output"
	"k8s.io/klog"
)

var alreadyInstalled = regexp.MustCompile("^Error: (.*?) already installed")

type Handler struct {
	Config   configuration
	Packages []Package
}

type configuration struct {
	// Defines a list of additional taps to install
	Taps               Taps `json:"taps;omitempty"`
	PrintCommandOutput bool `json:"printCommandOutput"`
}

type Package struct {
	Name    string `json:"name"`
	Tap     string `json:"tap,omitempty"`
	Version string `json:"version,omitempty"`
}

// format for printing packages is
// [tap]/<name>@[version]
func (p Package) String() string {
	var s string
	if p.Tap != "" {
		s += p.Tap + "/"
	}
	s += p.Name
	if p.Version != "" {
		s += "@" + p.Version
	}
	return s
}

// Init initialises the handler. If the packagelist should be nil, it adds itself
// to that list
func (b *Handler) Init(config *json.RawMessage, packages *json.RawMessage) error {
	if config != nil {
		err := json.Unmarshal([]byte(*config), &b.Config)
		if err != nil {
			return errors.Wrapf(err, "could not parse config %s", config)
		}
		klog.V(4).Infof("Brew: unmarshaled config %#v", b.Config)
	}

	if packages != nil {
		err := json.Unmarshal([]byte(*packages), &b.Packages)
		if err != nil {
			return errors.Wrapf(err, "could not parse packages %s", packages)
		}
		klog.V(4).Infof("Brew: Added package list %v", b.Packages)
	}

	err := b.Config.Taps.sync()
	return err
}

// New returns a handler with the default settings. They will be overwritten
// (if set) on Init()
func New() *Handler {
	return &Handler{
		// the default handler config
		Config: configuration{},
	}
}

// Install the packages and add them to the index. If an error occurs while installing
// a package, the other packages will be handled / installed nonetheless
func (b *Handler) Install(pkgs ...string) (packageList *json.RawMessage, err error) {
	return b.do(b.install, b.addToIndex, pkgs...)
}

// Remove a package from the system by parsing the given name and finding out the
// binary name. then remove the package from the index
func (b *Handler) Remove(pkgs ...string) (packageList *json.RawMessage, err error) {
	return b.do(b.remove, b.removeFromIndex, pkgs...)
}

// Upgrade a package, if it is in the index. Returns an error if a package
// should not exist in the index, but still processes all other packages
func (b *Handler) Upgrade(pkgs ...string) (packageList *json.RawMessage, err error) {
	return b.do(b.upgrade, b.upgradeIndex, pkgs...)
}

// do the package action and indexAction for a list of packages, handling
// errors and marshaling the index in the end
func (b *Handler) do(packageAction func(Package) error, indexAction func(Package), pkgs ...string) (*json.RawMessage, error) {
	var pError collection.Error
	packages, err := b.getPackages(pkgs...)
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
	raw, err := json.Marshal(b.Packages)
	if err != nil {
		return nil, err
	}
	msg := json.RawMessage(raw)
	return &msg, pError.IfNotEmpty()
}

func (b *Handler) install(pkg Package) error {
	output.Info("📦 Brew\t\tInstalling Package %s", pkg)
	p := pkg.Name
	if pkg.Version != "" {
		p += "@" + pkg.Version
	}
	if pkg.Tap != "" {
		p = pkg.Tap + "/" + p
	}

	c := []string{"brew", "install", p}
	out, err := cmd.Execute(c, cmd.DirectPrint(bool(klog.V(5)) || b.Config.PrintCommandOutput))
	// don't print the output twice if we have verbosity
	if err != nil && !b.Config.PrintCommandOutput && !bool(klog.V(5)) {
		output.Warn(out)
	}
	if err == nil {
		output.Success("📦 Brew\t\tInstalled Package %s", pkg)
	}
	if pkg.Version == "" {
		return err
	}

	// pin package if version is defined
	out, err = cmd.Execute([]string{"brew", "pin", p}, cmd.DirectPrint(bool(klog.V(5)) || b.Config.PrintCommandOutput))
	if err != nil && !b.Config.PrintCommandOutput && !bool(klog.V(5)) {
		output.Warn(out)
	}
	if err == nil {
		output.Success("📦 Brew\t\tPinned Package %s", pkg)
	}
	return err
}

func (b *Handler) remove(pkg Package) error {
	output.Info("📦 Brew\t\tRemoving Package %s", pkg)
	p := pkg.Name
	if pkg.Tap != "" {
		p = pkg.Tap + "/" + p
	}

	c := []string{"brew", "remove", p}
	out, err := cmd.Execute(c, cmd.DirectPrint(bool(klog.V(5)) || b.Config.PrintCommandOutput))
	// don't print the output twice if we have verbosity
	if err != nil && !b.Config.PrintCommandOutput && !bool(klog.V(5)) {
		output.Warn(out)
	}
	if err == nil {
		output.Success("📦 Brew\t\tInstalled Package %s", pkg)
	}
	return err
}

// returns the version of the package as defined in the index
func (b *Handler) IndexVersion(pkg Package) string {
	for _, p := range b.Packages {
		if p.Name == pkg.Name &&
			p.Tap == pkg.Tap {
			return p.Version
		}
	}
	return ""
}

func (b *Handler) upgrade(pkg Package) error {
	definedVersion := b.IndexVersion(pkg)
	if definedVersion == pkg.Version && pkg.Version != "" {
		output.Warn("📦 Brew\t\tNot upgrading package because it is pinned: %s", pkg)
		return nil
	}

	output.Info("📦 Brew\t\tUpgrading package %s", pkg)
	p := pkg.Name
	if pkg.Version != "" {
		p += "@" + pkg.Version
	}
	if pkg.Tap != "" {
		p = pkg.Tap + "/" + p
	}

	c := []string{"brew", "upgrade", p}
	out, err := cmd.Execute(c, cmd.DirectPrint(bool(klog.V(5)) || b.Config.PrintCommandOutput))
	if alreadyInstalled.MatchString(out) {
		output.Success("📦 Brew\t\tNo upgrade was needed for package %s", pkg)
		return nil
	}
	// don't print the output twice if we have verbosity
	if err != nil && !b.Config.PrintCommandOutput && !bool(klog.V(5)) {
		output.Warn(out)
	}
	if err == nil {
		output.Success("📦 Brew\t\tUpgraded Package %s", pkg)
	}
	return err
}

func (b *Handler) upgradeIndex(p Package) {
	for _, pkg := range b.Packages {
		if pkg.Name == p.Name &&
			pkg.Tap == p.Tap {
			b.addToIndex(p)
		}
	}
}

func (b *Handler) removeFromIndex(p Package) {
	klog.V(6).Infof("Brew: Removing Package %s from index", p)
	for idx, pkg := range b.Packages {
		if pkg == p {
			b.Packages = append(b.Packages[:idx], b.Packages[idx+1:]...)
			klog.V(6).Infof("Brew: Removed Package %s from index", p)
			return
		}
	}
	klog.V(6).Infof("Brew: Package %s has not been defined in index", p)
}

func (b *Handler) addToIndex(p Package) {
	klog.V(6).Infof("Brew: Adding package %s to index", p)
	for idx, pkg := range b.Packages {
		if pkg.Name == p.Name && pkg.Tap == p.Tap {
			b.Packages[idx] = p
			klog.V(6).Infof("Brew: Package %s already in index, overwriting old version %v", p, pkg.Version)
			return
		}
	}
	b.Packages = append(b.Packages, p)
}

// getPackages returns either a parsed list of the given packages
// or all the packages defined in the index, if no argument is given.
// no matter what is returned, the slice is save to use and modify.
// it can return an error for some packages, but always tries to parse
// all packages. All successfully parsed packages will be added to the
// slice, all others will be described in the error
func (b *Handler) getPackages(pkgs ...string) ([]Package, error) {
	var e collection.Error
	if len(pkgs) == 0 {
		klog.V(5).Infof("Brew: No packages defined, using all packages from go handler")
		// make it safe to modify
		packages := make([]Package, len(b.Packages))
		copy(packages, b.Packages)
		return packages, nil
	}

	var packages []Package
	for _, pkg := range pkgs {
		p, err := parse(pkg)
		if err != nil {
			klog.V(6).Infof("Brew: Error when parsing %v: %v", pkg, err)
			e.Add(pkg, err)
			continue
		}
		klog.V(6).Infof("Brew: Successfully parsed package %v (%s)", pkg, p)
		packages = append(packages, p)
	}
	return packages, e.IfNotEmpty()
}

func parse(pkg string) (Package, error) {
	var p Package
	if strings.Contains(pkg, "@") {
		v := strings.Split(pkg, "@")
		if len(v) > 2 {
			return p, errors.Errorf("invalid format for a package, too many '@': %v", pkg)
		}
		p.Version = v[1]
		pkg = v[0]
	}

	if strings.Contains(pkg, "/") {
		t := strings.Split(pkg, "/")
		p.Tap = strings.Join(t[:len(t)-1], "/")
		pkg = t[len(t)-1]
	}
	p.Name = pkg
	return p, nil
}
