package brew

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/pkg/output"
	"k8s.io/klog"
)

var alreadyInstalled = regexp.MustCompile("(.*?)Error: (.*?) already installed")

type Formula struct {
	Name    string `json:"name"`
	Tap     string `json:"tap,omitempty"`
	Version string `json:"version,omitempty"`
}

type Formulae []Formula

// format for printing packages is
// [tap]/<name>@[version]
func (f Formula) String() string {
	var s string
	if f.Tap != "" {
		s += f.Tap + "/"
	}
	s += f.Name
	if f.Version != "" {
		s += "@" + f.Version
	}
	return s
}

func (f Formula) fullname() string {
	if f.Tap != "" {
		return f.Tap + "/" + f.Name
	}
	return f.Name
}

func (f Formula) unpin() error {
	_, err := cmd.Execute([]string{"brew", "unpin", f.fullname()})
	return errors.Wrapf(err, "could not unpin formula %s", f)
}

func (f Formula) pin() error {
	_, err := cmd.Execute([]string{"brew", "pin", f.fullname()})
	return errors.Wrapf(err, "could not pin formula %s", f)
}

func (f Formula) install(printOutput bool) error {
	_, err := brewExec("install", f.String(), printOutput)
	return errors.Wrapf(err, "could not install formula %s", f)
}

func (f Formula) remove(printOutput bool) error {
	_, err := brewExec("remove", f.fullname(), printOutput)
	return errors.Wrapf(err, "could not remove formula %s", f)
}

func (f Formula) upgrade(printOutput bool) error {
	// code from brewExec, but with additional error handling
	out, err := cmd.Execute(
		[]string{"brew", "upgrade", f.String()},
		cmd.DirectPrint(bool(klog.V(5)) || printOutput),
	)
	// only print output if error occured and we have
	// not printed the output already
	if err != nil {
		if alreadyInstalled.MatchString(out) {
			return ErrNoUpgradeNeeded
		}
		if !printOutput && !bool(klog.V(5)) {
			output.Warn(out)
		}
	}
	return errors.Wrapf(err, "could not upgrade formula %s", f)
}

func brewExec(action, formula string, printOutput bool) (out string, err error) {
	out, err = cmd.Execute(
		[]string{"brew", action, formula},
		cmd.DirectPrint(bool(klog.V(5)) || printOutput),
	)
	// only print output if error occured and we have
	// not printed the output already
	if err != nil && !printOutput && !bool(klog.V(5)) {
		output.Warn(out)
	}
	return out, errors.Wrapf(err, "could not %s %s", action, formula)
}

func parse(formula string) (Formula, error) {
	var f Formula
	// extract the version
	if strings.Contains(formula, "@") {
		v := strings.Split(formula, "@")
		if len(v) > 2 {
			return f, errors.Errorf("invalid format for a package, too many '@': %v", formula)
		}
		formula = v[0]
		f.Version = v[1]
	}

	// extract the tap
	if strings.Contains(formula, "/") {
		t := strings.Split(formula, "/")
		f.Tap = strings.Join(t[:len(t)-1], "/")
		formula = t[len(t)-1]
	}
	f.Name = formula
	return f, nil
}
