package brew

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
)

var defaultTaps = []string{"homebrew/core"}

// Tap contains of a name, and if it should be
// cloned fully or shallow (see `brew tap -h` for
// further info)
type Tap struct {
	Name string `json:"name;omitempty"`
	Full bool   `json:"full"`
}

func (t Tap) String() string {
	return t.Name
}

type Taps []Tap

func (t Taps) names() (names []string) {
	for _, tap := range t {
		names = append(names, tap.Name)
	}
	return names
}

// getInstalledTaps returns all installed taps, EXCEPT
// the default taps
func getInstalledTaps() (taps []string, err error) {
	list, err := cmd.Execute([]string{"brew", "tap"})
	if err != nil {
		return taps, errors.Wrapf(err, "output: %v", list)
	}
	// remove defaultTaps
	for _, dt := range defaultTaps {
		list = strings.ReplaceAll(list, dt+"\n", "")
	}
	taps = strings.Split(list, "\n")
	// remove any empty strings or newlines
	var rt []string
	for _, t := range taps {
		if t != "" && t != "\n" {
			rt = append(rt, t)
		}
	}
	return rt, nil
}

func (t Tap) install() error {
	c := []string{"brew", "tap", t.Name}
	if t.Full {
		c = append(c, "--full")
	}
	_, err := cmd.Execute(c)
	return errors.Wrapf(err, "could not install tap %s", t)
}

func (t Tap) remove() error {
	_, err := cmd.Execute([]string{"brew", "untap", t.Name})
	return errors.Wrapf(err, "could not remove tap %s", t)
}

// sync taps, meaning install taps that are defined in
// Taps but not installed, and remove taps that are installed
// but not defined in Taps
func (t Taps) sync() error {
	installedTaps, err := getInstalledTaps()
	if err != nil {
		return errors.Wrap(err, "could not get list of installed taps")
	}

	missing, spare := filterTaps(installedTaps, t.names())

	for _, m := range missing {
		if err := t.Tap(m).install(); err != nil {
			return errors.Wrap(err, "could not install missing tap")
		}
	}

	for _, s := range spare {
		err := Tap{Name: s}.remove()
		if err != nil {
			return errors.Wrap(err, "could not remove spare tap")
		}
	}
	return nil
}

// Tap returns the tap as defined in Taps, or,
// if it should not exist, a new tap
func (t Taps) Tap(name string) Tap {
	for _, tap := range t {
		if tap.Name == name {
			return tap
		}
	}
	return Tap{Name: name}
}

func filterTaps(installedTaps, desiredTaps []string) (missingTaps, spareTaps []string) {
	for _, inst := range installedTaps {
		if !contains(inst, desiredTaps) {
			spareTaps = append(spareTaps, inst)
		}
	}
	for _, des := range desiredTaps {
		if !contains(des, installedTaps) {
			missingTaps = append(missingTaps, des)
		}
	}

	return missingTaps, spareTaps
}

func contains(s string, ls []string) bool {
	for _, elem := range ls {
		if elem == s {
			return true
		}
	}
	return false
}
