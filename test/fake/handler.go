// Package fake contains fake implementations of resources used in this
// program.
package fake

import (
	"encoding/json"
	"errors"
)

type Handler struct {
	Config
	Packages []Package
}

type Config struct {
	WorkingDir string `json:"workingDir"`
}
type Package struct {
	Name string `json:"url"`
}

var DefaultSettings = Config{"tmp"}
var DefaultPackages = []Package{{"fakePackage1"}, {"fakePackage2"}}

// DefaultSettingsRaw for this handler, ready to use
var DefaultSettingsRaw = func() *json.RawMessage {
	d, _ := json.Marshal(DefaultSettings)
	msg := json.RawMessage(d)
	return &msg
}()

var DefaultPackagesRaw = func() *json.RawMessage {
	d, _ := json.Marshal(DefaultPackages)
	msg := json.RawMessage(d)
	return &msg
}()
var EmptyPackages = func() *json.RawMessage {
	msg := json.RawMessage(`[]`)
	return &msg
}()

func (h *Handler) Init(config *json.RawMessage, packages *json.RawMessage) error {
	if config != nil {
		err := json.Unmarshal([]byte(*config), &h.Config)
		if err != nil {
			return err
		}
	}
	if packages == nil {
		return errors.New("no packages declared, should not get empty package list")
	}
	return json.Unmarshal([]byte(*packages), &h.Packages)
}

func (h *Handler) Install(pkgs ...string) (*json.RawMessage, error) {
	for _, pkg := range pkgs {
		h.Packages = append(h.Packages, Package{pkg})
	}
	return h.marshalPackages()
}

// Remove fails on the first package that has not been found and
// does not process all packages on failure!
func (h *Handler) Remove(pkgs ...string) (*json.RawMessage, error) {
	for _, pkg := range pkgs {
		removed := false
		for i, p := range h.Packages {
			if pkg == p.Name {
				h.Packages = append(h.Packages[:i], h.Packages[i+1:]...)
				removed = true
				break
			}
		}
		if !removed {
			return nil, errors.New("no package found in index to remove")
		}
	}

	return h.marshalPackages()
}

// Upgrade adds a "+" at the end of the package string
func (h *Handler) Upgrade(pkgs ...string) (*json.RawMessage, error) {
	if len(pkgs) == 0 {
		return h.upgradeAll()
	}
	for _, pkg := range pkgs {
		for i := range h.Packages {
			if h.Packages[i].Name == pkg {
				h.Packages[i].Name += "+"
			}
		}
	}
	return h.marshalPackages()
}

func (h *Handler) upgradeAll() (*json.RawMessage, error) {
	for i := range h.Packages {
		h.Packages[i].Name += "+"
	}
	return h.marshalPackages()
}

func (h *Handler) marshalPackages() (*json.RawMessage, error) {
	pkgList, err := json.Marshal(h.Packages)
	if err != nil {
		return nil, err
	}
	rM := json.RawMessage(pkgList)
	return &rM, nil

}
