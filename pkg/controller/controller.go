package controller

import (
	"encoding/json"

	"github.com/tommyknows/packa/pkg/collection"
	"github.com/tommyknows/packa/pkg/output"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

// Controller serves as a unifying interface of different
// "handlers", implementations of different package managers.
type Controller struct {
	configuration *Configuration
	handlers      map[string]*handler
}

type handler struct {
	PackageHandler
	// if the init function of the handler has been run
	initialised bool
}

func (h *handler) setInitialised() {
	h.initialised = true
}

// PackageType is an interface that defines the methods for package types.
// for every operation, the following statements are true:
//   - If pkg is an empty string, install all packages as defined
//     in the handler's package list.
//   - As long as the returned packageList is not nil, it will be added
//     back to the config.
//   - The given packages will most likely need to be parsed (e.g. separating
//     the package name and version)
type PackageHandler interface {
	// Init the package handler with the given config and packages.
	Init(config *json.RawMessage, packages *json.RawMessage) error
	// Install the given packages on the system.
	// See docs on PackageHandler
	Install(pkgs ...string) (packageList *json.RawMessage, err error)
	// Remove the given packages on the system.
	// See docs on PackageHandler
	Remove(pkgs ...string) (packageList *json.RawMessage, err error)
	// Upgrade the given packages on the system.
	// See docs on PackageHandler
	Upgrade(pkgs ...string) (packageList *json.RawMessage, err error)
}

// handlerOperation is a function taken from the PackageHandler interface.
type handlerOperation func(handler PackageHandler, pkgs ...string) (packageList *json.RawMessage, err error)

// functional-style options
type option func(*Controller) error

// New creates a new Controller, ready to use. Uses functional options
// to set different fields
func New(opts ...option) (*Controller, error) {
	ctl := &Controller{
		configuration: defaultConfig(),
		handlers:      make(map[string]*handler),
	}
	for _, opt := range opts {
		err := opt(ctl)
		if err != nil {
			return nil, errors.Wrapf(err, "could not add option")
		}
	}
	return ctl, nil
}

// RegisterHandlers registers the given handlers on the controller
func RegisterHandlers(handlers map[string]PackageHandler) option {
	return func(ctl *Controller) error {
		for name, h := range handlers {
			klog.V(2).Infof("Registering handler %v", name)
			ctl.handlers[name] = &handler{h, false}
		}
		return nil
	}
}

// Close contains cleanup tasks that should be done when the command ends.
// right now, this is mainly saving the config state to file
func (ctl *Controller) Close() error {
	klog.V(3).Infof("Closing controller")
	err := ctl.configuration.save()
	return errors.Wrapf(err, "could not save config")
}

// PrintPackages of the specified handlers. If no handler is specified,
// print all packages from all handlers
func (ctl *Controller) PrintPackages(handlers ...string) error {
	if len(handlers) == 0 {
		klog.V(1).Infof("Printing packages of all handlers")
		for h := range ctl.handlers {
			handlers = append(handlers, h)
		}
	}

	var cerr collection.Error
	for _, h := range handlers {
		klog.V(2).Infof("Printing packages of handler %v", h)
		pkgs, ok := ctl.configuration.Packages[h]
		if !ok {
			cerr.Add(h, errors.New("handler does not exist"))
			continue
		}
		if pkgs == nil {
			output.Info("Handler %v does not specify any packages", h)
			continue
		}

		out, err := yaml.JSONToYAML([]byte(*pkgs))
		if err != nil {
			return err
		}
		output.Info("%s:\n%s", h, out)
	}
	return cerr.IfNotEmpty()
}

// Install the package with the handler
// If pkg is empty string, install all packages that are defined in
// the handler's package list
func (ctl *Controller) Install(handler string, pkgs ...string) error {
	klog.V(2).Infof("Installing package(s) %v on handler %v", pkgs, handler)
	return ctl.handlerDo(PackageHandler.Install, handler, pkgs...)
}

// Remove the package with the handler
// If pkg is empty string, upgrade all packages that are defined in
// the handler's package list
func (ctl *Controller) Remove(handler string, pkgs ...string) error {
	klog.V(2).Infof("Removing packages %v on handler %v", pkgs, handler)
	return ctl.handlerDo(PackageHandler.Remove, handler, pkgs...)
}

// Upgrade the given package with the handler.
// If pkg is empty string, upgrade all packages that are defined in
// the handler's package list.
func (ctl *Controller) Upgrade(handler string, pkgs ...string) error {
	klog.V(2).Infof("Upgrading packages %v on handler %v", pkgs, handler)
	return ctl.handlerDo(PackageHandler.Upgrade, handler, pkgs...)
}

// UpgradeAll upgrades all packages from all handlers.
func (ctl *Controller) UpgradeAll() error {
	klog.V(2).Infof("Upgrading all packages")
	var ce collection.Error
	for name, _ := range ctl.handlers {
		err := ctl.handlerDo(PackageHandler.Upgrade, name)
		if err != nil {
			ce.Add(name, err)
		}
	}

	return ce.IfNotEmpty()
}

// handlerDo takes operations defined on the handlerInterface and executes
// them accordingly. Does all necessary safetychecks and config-modifications.
func (ctl *Controller) handlerDo(f handlerOperation, handler string, pkgs ...string) error {
	// check if the handler even exists / got registered
	if ctl.handlers[handler] == nil {
		return errors.Errorf("handler \"%v\" does not exist or has not been registered", handler)
	}

	// initialise the handler if it has not been initialised
	if !ctl.handlers[handler].initialised {
		klog.V(2).Infof("Initialising handler %v", handler)
		err := ctl.initialiseHandler(handler)
		if err != nil {
			return err
		}
	}

	// execute the actual function and update the index
	pkgList, err := f(ctl.handlers[handler], pkgs...)
	if pkgList != nil {
		klog.V(3).Infof("Appending new packagelist to handler %v", handler)
		ctl.configuration.Packages[handler] = pkgList
	}
	return errors.Wrapf(err, "error executing action on handler %v", handler)
}

// initialiseHandler initialises the handler with the given name,
// calling its Init method with the settings and packages as defined
// in the configuration.
func (ctl *Controller) initialiseHandler(name string) error {
	settings := ctl.configuration.Settings.Handler[name]
	packages := ctl.configuration.Packages[name]

	err := ctl.handlers[name].Init(settings, packages)
	if err != nil {
		return errors.Wrapf(err, "could not initialise handler %v", name)
	}
	ctl.handlers[name].setInitialised()
	klog.V(2).Infof("Handler %v successfully initialised", name)
	return nil
}
