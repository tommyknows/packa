package brew

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
)

var defaultTaps = []string{"homebrew/core"}

// tap contains of a name, and if it should be
// cloned fully or shallow (see `brew tap -h` for
// further info)
type tap struct {
	Name string `json:"name;omitempty"`
	Full bool   `json:"full"`
}

func (t tap) String() string {
	return t.Name
}

type taps []tap

func (t taps) names() (names []string) {
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

func (t tap) install() error {
	c := []string{"brew", "tap", t.Name}
	if t.Full {
		c = append(c, "--full")
	}
	_, err := cmd.Execute(c)
	return errors.Wrapf(err, "could not install tap %s", t)
}

func (t tap) remove() error {
	_, err := cmd.Execute([]string{"brew", "untap", t.Name})
	return errors.Wrapf(err, "could not remove tap %s", t)
}

// sync taps, meaning install taps that are defined in
// Taps but not installed, and remove taps that are installed
// but not defined in Taps
func (t taps) sync() error {
	installedTaps, err := getInstalledTaps()
	if err != nil {
		return errors.Wrap(err, "could not get list of installed taps")
	}

	missing, spare := filterTaps(installedTaps, t.names())

	for _, m := range missing {
		if err := t.tap(m).install(); err != nil {
			return errors.Wrap(err, "could not install missing tap")
		}
	}

	for _, s := range spare {
		err := tap{Name: s}.remove()
		if err != nil {
			return errors.Wrap(err, "could not remove spare tap")
		}
	}
	return nil
}

// tap returns the tap as defined in Taps, or,
// if it should not exist, a new tap
func (t taps) tap(name string) tap {
	for _, tap := range t {
		if tap.Name == name {
			return tap
		}
	}
	return tap{Name: name}
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
