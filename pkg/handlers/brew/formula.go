package brew

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/pkg/output"
	"k8s.io/klog"
)

var alreadyInstalled = regexp.MustCompile("^Error: (.*?) already installed")

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
	return errors.Wrapf(err, "could not unpin package %v", f)
}

func (f Formula) pin() error {
	_, err := cmd.Execute([]string{"brew", "pin", f.fullname()})
	return errors.Wrapf(err, "could not pin package %v", f)
}

func (f Formula) install(printOutput bool) error {
	_, err := brewExec("install", f.String(), printOutput)
	return err
}

func (f Formula) remove(printOutput bool) error {
	_, err := brewExec("remove", f.fullname(), printOutput)
	return err
}

func (f Formula) upgrade(printOutput bool) error {
	out, err := brewExec("upgrade", f.String(), printOutput)
	if alreadyInstalled.MatchString(out) {
		return ErrNoUpgradeNeeded
	}
	return err
}

func brewExec(action, pkg string, printOutput bool) (out string, err error) {
	out, err = cmd.Execute(
		[]string{"brew", action, pkg},
		cmd.DirectPrint(bool(klog.V(5)) || printOutput),
	)
	// don't print the output twice if we have verbosity
	if err != nil && !printOutput && !bool(klog.V(5)) {
		output.Warn(out)
	}
	return out, errors.Wrapf(err, "could not %s %s", action, pkg)
}

func parse(pkg string) (Formula, error) {
	var f Formula
	if strings.Contains(pkg, "@") {
		v := strings.Split(pkg, "@")
		if len(v) > 2 {
			return f, errors.Errorf("invalid format for a package, too many '@': %v", pkg)
		}
		f.Version = v[1]
		pkg = v[0]
	}

	if strings.Contains(pkg, "/") {
		t := strings.Split(pkg, "/")
		f.Tap = strings.Join(t[:len(t)-1], "/")
		pkg = t[len(t)-1]
	}
	f.Name = pkg
	return f, nil
}
