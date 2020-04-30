package brew

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/pkg/collection"
	"github.com/tommyknows/packa/pkg/output"
	"k8s.io/klog"
)

type errorConst string

const ErrNoUpgradeNeeded errorConst = "no upgrade was needed"

func (e errorConst) Error() string {
	return string(e)
}

type Handler struct {
	Config   configuration
	Formulae formulae
}

type configuration struct {
	// Defines a list of additional taps to install
	Taps               taps `json:"taps;omitempty"`
	PrintCommandOutput bool `json:"printCommandOutput;omitempty"`
	UpdateOnInit       bool `json:"updateOnInit;omitempty"`
}

// Init initialises the handler.
func (b *Handler) Init(config *json.RawMessage, formulaeList *json.RawMessage) error {
	if config != nil {
		err := json.Unmarshal([]byte(*config), &b.Config)
		if err != nil {
			return errors.Wrapf(err, "could not parse config %s", config)
		}
		klog.V(4).Infof("Brew: unmarshaled config %#v", b.Config)
	}

	if formulaeList != nil {
		err := json.Unmarshal([]byte(*formulaeList), &b.Formulae)
		if err != nil {
			return errors.Wrapf(err, "could not parse formulae %s", formulaeList)
		}
		klog.V(4).Infof("Brew: Added formulae list %v", b.Formulae)
	}

	if b.Config.UpdateOnInit {
		if err := updateBrew(); err != nil {
			return errors.Wrapf(err, "auto-update failed")
		}
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

// Install the formulae and add them to the index. If an error occurs while installing
// a formula, the other formulae will be handled / installed nonetheless
func (b *Handler) Install(pkgs ...string) (formulaList *json.RawMessage, err error) {
	return b.do(b.install, b.addToIndex, pkgs...)
}

// Remove formulae from the system. If an error occurs while installing
// a formula, the other formulae will be handled / installed nonetheless
func (b *Handler) Remove(pkgs ...string) (formulaList *json.RawMessage, err error) {
	return b.do(b.remove, b.removeFromIndex, pkgs...)
}

// Upgrade a formula, if it is in the index. Returns an error if a formula
// should not exist in the index, but still processes all other formulae
func (b *Handler) Upgrade(pkgs ...string) (formulaList *json.RawMessage, err error) {
	return b.do(b.upgrade, b.upgradeIndex, pkgs...)
}

// do the formula action and indexAction for a list of formulae, handling
// errors and marshaling the index in the end
// NOTE: this code is more or less exactly the same to the method in the goget-package...
func (b *Handler) do(formulaAction func(formula) error, indexAction func(formula), pkgs ...string) (*json.RawMessage, error) {
	var pError collection.Error
	forms, err := b.getFormulae(pkgs...)
	if err != nil {
		pe, ok := err.(*collection.Error)
		if !ok {
			// we always expect a collection.Error from getFormulae
			return nil, errors.Wrapf(err, "unexpected error occured when getting formulae list")
		}
		pError.Merge(*pe)
	}

	// should never occur, but as a safety check
	if formulaAction == nil {
		return nil, errors.New("no formula action defined")
	}

	for _, p := range forms {
		// execute the formulaAction and then the index action, if applicable
		switch err := formulaAction(p); {
		case indexAction == nil:
			klog.V(5).Infof("Brew: Not executing index action because none is defined")
		case err != nil:
			klog.V(4).Infof("Brew: Error while executing formula action for %v, adding error to collection", p.String())
			pError.Add(p.String(), err)
			klog.V(5).Infof("Brew: Not executing index action because of error on formula action")
		default:
			indexAction(p)
		}
	}

	klog.V(6).Infof("Brew: Marshaling formulae")
	raw, err := json.Marshal(b.Formulae)
	if err != nil {
		return nil, err
	}
	msg := json.RawMessage(raw)
	return &msg, pError.IfNotEmpty()
}

func (b *Handler) install(f formula) error {
	output.Info("📦 Brew\t\tInstalling formula %s", f)
	err := f.install(b.Config.PrintCommandOutput)
	if err != nil {
		return err
	}
	output.Success("📦 Brew\t\tInstalled formula %s", f)

	if f.Version == "" {
		return err
	}

	// pin package if version is defined
	err = f.pin()
	if err == nil {
		output.Success("📦 Brew\t\tPinned formula %s", f)
	}
	return err
}

func (b *Handler) remove(f formula) error {
	output.Info("📦 Brew\t\tRemoving formula %s", f)
	err := f.uninstall(b.Config.PrintCommandOutput)
	if err == nil {
		output.Success("📦 Brew\t\tRemoved formula %s", f)
	}
	return err
}

func (b *Handler) upgrade(f formula) error {
	if definedVersion := b.indexVersion(f); definedVersion != "" {
		if f.Version == definedVersion {
			output.Warn("📦 Brew\t\tNot upgrading package because it is pinned: %s", f)
			return nil
		}

		output.Warn("📦 Brew\t\tUnpinning pinned package %s for upgrade", f)
		err := f.unpin()
		if err != nil {
			return errors.Wrapf(err, "could not unpin package %v", f.Name)
		}
	}

	output.Info("📦 Brew\t\tUpgrading package %s", f)
	err := f.upgrade(b.Config.PrintCommandOutput)
	if err != nil && err != ErrNoUpgradeNeeded {
		return err
	}
	if err == ErrNoUpgradeNeeded {
		output.Success("📦 Brew\t\tNo upgrade was needed for formula %s", f)
		err = nil
	} else {
		output.Success("📦 Brew\t\tUpgraded Package %s", f)
	}

	if f.Version == "" {
		return err
	}

	// pin package if version is defined
	err = f.pin()
	if err == nil {
		output.Success("📦 Brew\t\tPinned Package %s", f)
	}
	return err
}

// returns the version of the package as defined in the index
func (b *Handler) indexVersion(f formula) string {
	for _, form := range b.Formulae {
		if form.Name == f.Name &&
			form.Tap == f.Tap {
			return form.Version
		}
	}
	return ""
}

func (b *Handler) upgradeIndex(f formula) {
	for _, formula := range b.Formulae {
		if formula.Name == f.Name &&
			formula.Tap == f.Tap {
			b.addToIndex(f)
		}
	}
}

func (b *Handler) removeFromIndex(f formula) {
	klog.V(6).Infof("Brew: Removing Package %s from index", f)
	for idx, pkg := range b.Formulae {
		if pkg == f {
			b.Formulae = append(b.Formulae[:idx], b.Formulae[idx+1:]...)
			klog.V(6).Infof("Brew: Removed Package %s from index", f)
			return
		}
	}
	klog.V(6).Infof("Brew: Package %s has not been defined in index", f)
}

func (b *Handler) addToIndex(f formula) {
	klog.V(6).Infof("Brew: Adding package %s to index", f)
	for idx, pkg := range b.Formulae {
		if pkg.Name == f.Name && pkg.Tap == f.Tap {
			b.Formulae[idx] = f
			klog.V(6).Infof("Brew: Package %s already in index, overwriting old version %v", f, pkg.Version)
			return
		}
	}
	b.Formulae = append(b.Formulae, f)
}

// getFormulae returns either a parsed list of the given packages
// or all the packages defined in the index, if no argument is given.
// no matter what is returned, the slice is save to use and modify.
// it can return an error for some packages, but always tries to parse
// all packages. All successfully parsed packages will be added to the
// slice, all others will be described in the error
func (b *Handler) getFormulae(forms ...string) (formulae, error) {
	var e collection.Error
	if len(forms) == 0 {
		klog.V(5).Infof("Brew: No packages defined, using all packages from brew handler")
		// make it safe to modify
		packages := make(formulae, len(b.Formulae))
		copy(packages, b.Formulae)
		return packages, nil
	}

	var f formulae
	var isCask bool
	if forms[0] == "cask" {
		isCask = true
		forms = forms[1:]
	}

	for _, pkg := range forms {
		p, err := parse(pkg, isCask)
		if err != nil {
			klog.V(6).Infof("Brew: Error when parsing %v: %v", pkg, err)
			e.Add(pkg, err)
			continue
		}
		klog.V(6).Infof("Brew: Successfully parsed package %v (%s)", pkg, p)
		f = append(f, p)
	}
	return f, e.IfNotEmpty()
}

func updateBrew() error {
	_, err := cmd.Execute([]string{"brew", "update"})
	return errors.Wrapf(err, "could not update brew")
}
