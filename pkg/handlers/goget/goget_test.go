package goget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tommyknows/packa/pkg/collection"
	"github.com/tommyknows/packa/pkg/defaults"
)

func newTestHandler(t *testing.T) *Handler {
	is := is.New(t)

	h := New()
	tmpDir, err := ioutil.TempDir("", "")
	is.NoErr(err)
	cfg := json.RawMessage(`{"workingDir": "` + tmpDir + `"}`)
	err = h.Init(&cfg, nil)
	is.NoErr(err)
	return h
}

func TestInit(t *testing.T) {
	is := is.New(t)

	tests := map[string]struct {
		configRaw   json.RawMessage
		packagesRaw json.RawMessage
		config      configuration
		packages    []Package
		isErr       bool
	}{
		"empty": {
			configRaw:   json.RawMessage(`{}`),
			config:      configuration{defaults.WorkingDir(), false, false}, // default value from New() call
			packagesRaw: json.RawMessage(`[]`),
			packages:    []Package{},
			isErr:       false,
		},
		"nil-config": {
			configRaw:   json.RawMessage(``),
			config:      configuration{defaults.WorkingDir(), false, false}, // default value from New() call
			packagesRaw: json.RawMessage(`[]`),
			packages:    []Package{},
			isErr:       false,
		},
		"nil-packages": {
			configRaw:   json.RawMessage(`{}`),
			config:      configuration{defaults.WorkingDir(), false, false}, // default value from New() call
			packagesRaw: json.RawMessage(``),
			packages:    []Package{{"github.com/tommyknows/packa", "latest"}},
			isErr:       false,
		},
		"configed": {
			configRaw:   json.RawMessage(`{"workingDir": "/test"}`),
			config:      configuration{"/test", false, false},
			packagesRaw: json.RawMessage(`[]`),
			packages:    []Package{},
			isErr:       false,
		},
		"defined package": {
			configRaw:   json.RawMessage(`{"workingDir": "/test"}`),
			config:      configuration{"/test", false, false},
			packagesRaw: json.RawMessage(`[{"url": "github.com/test/test", "version": "latest"}]`),
			packages:    []Package{{"github.com/test/test", "latest"}},
			isErr:       false,
		},
		"invalid config": {
			configRaw:   json.RawMessage(`{"workingDir" "/test"}`),
			config:      configuration{},
			packagesRaw: json.RawMessage(`[]`),
			packages:    []Package{},
			isErr:       true,
		},
		"invalid packages": {
			configRaw:   json.RawMessage(``),
			config:      configuration{defaults.WorkingDir(), false, false}, // default value from New() call
			packagesRaw: json.RawMessage(`["test":"bla"]`),
			packages:    []Package{},
			isErr:       true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := New()
			cr := &tt.configRaw
			pr := &tt.packagesRaw
			if len(*cr) == 0 {
				cr = nil
			}
			if len(*pr) == 0 {
				pr = nil
			}
			err := h.Init(cr, pr)
			if tt.isErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			is.Equal(tt.config, h.Config)
			is.Equal(tt.packages, h.Packages)
		})
	}
}

// TODO: this test should be refactored, I feel like this needs & could
// be simplified
func TestDo(t *testing.T) {
	is := is.New(t)

	pkgCh := make(chan Package)
	idxCh := make(chan Package)
	pkgAction := func(returnError bool) func(Package) error {
		return func(p Package) error {
			var ce collection.Error
			pkgCh <- p
			if returnError {
				ce.Add("test", fmt.Errorf("artificial error"))
			}
			return ce.IfNotEmpty()
		}
	}

	idxAction := func(p Package) {
		idxCh <- p
	}

	h := newTestHandler(t)
	rm, err := h.do(nil, nil)
	is.True(rm == nil)
	is.True(err != nil)

	go func() {
		rm, err = h.do(pkgAction(false), idxAction)
		is.True(rm != nil)
		is.NoErr(err)
		close(idxCh)
		close(pkgCh)
	}()

	run := true
	for run {
		select {
		case p, ok := <-pkgCh:
			if !ok {
				pkgCh = nil
				continue
			}
			is.Equal(h.Packages[0].URL, p.URL)
		case p, ok := <-idxCh:
			if !ok {
				idxCh = nil
				continue
			}
			is.Equal(h.Packages[0].URL, p.URL)
		default:
			if pkgCh == nil && idxCh == nil {
				run = false
			}
		}
	}

	pkgCh = make(chan Package)
	idxCh = make(chan Package)

	go func() {
		rm, err = h.do(pkgAction(true), idxAction)
		is.True(rm != nil)
		is.True(err != nil)
		close(idxCh)
		close(pkgCh)
	}()

	run = true
	for run {
		select {
		case p, ok := <-pkgCh:
			if !ok {
				pkgCh = nil
				continue
			}
			is.Equal(h.Packages[0].URL, p.URL)
		case p, ok := <-idxCh:
			if !ok {
				idxCh = nil
				continue
			}
			is.Equal(h.Packages[0].URL, p.URL)
		default:
			if pkgCh == nil && idxCh == nil {
				run = false
			}
		}
	}

	rm, err = h.do(pkgAction(false), idxAction, "test@bla@x@")
	is.True(rm != nil)
	is.True(err != nil)
	ce, ok := err.(*collection.Error)
	if !ok {
		t.Fatalf("did not get a collection error, got %v (%T)", ce, ce)
	}
	is.True(strings.Contains(ce.Error(), "test@bla@x@"))
}

func TestParse(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		in    string
		out   Package
		isErr bool
	}{
		{
			"github.com/test/test@master",
			Package{"github.com/test/test", "master"},
			false,
		},
		{
			"github.com/test/test",
			Package{"github.com/test/test", ""},
			false,
		},
		{
			"test@test@latest",
			Package{},
			true,
		},
	}

	for _, tt := range tests {
		p, err := parse(tt.in)
		is.Equal(tt.out, p)
		is.Equal(tt.isErr, err != nil)
	}
}

func TestGetPackages(t *testing.T) {
	is := is.New(t)

	h := newTestHandler(t)
	pkgs, err := h.getPackages()
	is.NoErr(err)
	is.Equal(h.Packages, pkgs)
	pkgs[0].URL = "somethingdifferent"
	is.True(h.Packages[0] != pkgs[0])

	packages := []string{"test@master", "test2", "test3@v3.0.0"}
	pkgs, err = h.getPackages(packages...)
	is.NoErr(err)
	is.Equal(packages[0], pkgs[0].String())
	is.Equal(packages[1], pkgs[1].String())
	is.Equal(packages[2], pkgs[2].String())

	packages = append(packages, "invalid@package@url")
	pkgs, err = h.getPackages(packages...)
	is.True(err != nil)
	is.Equal(packages[0], pkgs[0].String())
	is.Equal(packages[1], pkgs[1].String())
	is.Equal(packages[2], pkgs[2].String())
}

func TestIndexActions(t *testing.T) {
	is := is.New(t)

	h := newTestHandler(t)
	tp1 := Package{"test.com/test", "latest"}
	h.addToIndex(tp1)
	is.Equal(tp1, h.Packages[len(h.Packages)-1])
	tp2 := h.Packages[0]
	tp2.Version = "latest"
	h.addToIndex(tp2)
	is.Equal(tp2, h.Packages[0])

	tp2.Version = "newlatest"
	h.upgradeIndex(tp2)
	is.Equal(tp2, h.Packages[0])

	h.removeFromIndex(tp2)
	is.Equal(tp1, h.Packages[0])
	is.True(len(h.Packages) == 1)

	// package should not be added to the index
	// because we just deleted it
	h.upgradeIndex(tp2)
	is.Equal(tp1, h.Packages[0])
	is.True(len(h.Packages) == 1)
}

func TestHasActions(t *testing.T) {
	is := is.New(t)

	h := newTestHandler(t)
	p := h.Packages[0]
	is.True(h.hasURL(p))
	is.True(h.has(p))

	p.Version = "somethingelse"
	is.True(h.hasURL(p))
	is.True(h.has(p) == false)

	p.URL = "somethingelse"
	is.True(h.hasURL(p) == false)
	is.True(h.has(p) == false)
}

func TestMatchSemVer(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		in  string
		out bool
	}{
		{"v1.0.0", true},
		{"v0.0.12", true},
		{"v0.test", false},
		{"v0.0.0-commitid", true},
		{"branchname", false},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := matchSemVer(tt.in)
			is.Equal(tt.out, res)
		})
	}
}

func TestExtractBinaryName(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		in  string
		out string
	}{
		{"github.com/test/testbin", "testbin"},
		{"github.com/test/testbin/", "testbin"},
		{"github.com/test/testbin/v3", "testbin"},
		{"github.com/test/testbin/v3/", "testbin"},
		{"url.com/binary", "binary"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := extractBinaryName(tt.in)
			is.Equal(tt.out, res)
		})
	}
}
