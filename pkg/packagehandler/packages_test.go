package packages

import (
	"fmt"
	"testing"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"git.ramonruettimann.ml/ramon/packa/test/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewPackage(t *testing.T) {
	tests := []struct {
		url string
		pkg Package
	}{
		{
			url: "test.com/package@latest",
			pkg: Package{
				&config.Package{
					URL:     "test.com/package",
					Version: "latest",
				},
				nil,
			},
		},
		{
			url: "abc.def/another/subpackage@v0.0.1",
			pkg: Package{
				&config.Package{
					URL:     "abc.def/another/subpackage",
					Version: "v0.0.1",
				},
				nil,
			},
		},
		{
			url: "abc.def/nogiven/version",
			pkg: Package{
				&config.Package{
					URL:     "abc.def/nogiven/version",
					Version: "latest",
				},
				nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			p := NewPackage(tt.url, nil)
			assert.Equal(t, tt.pkg.URL, p.URL)
			assert.Equal(t, tt.pkg.Version, p.Version)
		})
	}
}

func TestInstallPackage(t *testing.T) {
	tests := []struct {
		pkg     Package
		version string
		err     error
	}{
		{
			pkg: Package{
				&config.Package{
					URL:     "test.com/no/bla",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("go: extracting test.com/no/bla v0.0.1\n", "v0.0.1", nil),
			},
			version: "v0.0.1",
			err:     nil,
		},
		{
			pkg: Package{
				&config.Package{
					URL:     "test.com/no/bla",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("multiline\ngo: extracting test.com/no/bla v0.0.1\n", "v0.0.1", nil),
			},
			version: "v0.0.1",
			err:     nil,
		},
		{
			pkg: Package{
				&config.Package{
					URL:     "extracttest ",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("some\nmultiline\noutput\n", "~v0.0.1", nil),
			},
			version: "~v0.0.1",
			err:     nil,
		},
		{
			pkg: Package{
				&config.Package{
					URL:     "error output",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("help: error\n", "", fmt.Errorf("exit 1")),
			},
			version: "",
			err:     errors.Wrapf(fmt.Errorf("help: error\nexit 1"), "could not install package error output"),
		},
		{
			pkg: Package{
				&config.Package{
					URL:     "someurl",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("", "v0.0.1", nil),
			},
			version: "v0.0.1",
			err:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.pkg.URL, func(t *testing.T) {
			err := tt.pkg.Install()
			assert.Equal(t, errors.Cause(tt.err), errors.Cause(err))
		})
	}
}

func TestRemovePackage(t *testing.T) {
	tests := []struct {
		pkg    Package
		binary string
		err    error
	}{
		{
			pkg: Package{
				&config.Package{
					URL:     "test.com/no/bla",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("go: extracting test.com/no/bla v0.0.1\n", "v0.0.1", nil),
			},
			binary: "bla",
			err:    nil,
		},
		{
			pkg: Package{
				&config.Package{
					URL:     "test.com/no/anotherbinary",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("multiline\ngo: extracting test.com/no/bla v0.0.1\n", "", nil),
			},
			binary: "anotherbinary",
			err:    nil,
		},
		{
			pkg: Package{
				&config.Package{
					URL:     "test.com/name",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("", "", fmt.Errorf("exit 1")),
			},
			binary: "name",
			err:    fmt.Errorf("exit 1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.pkg.URL, func(t *testing.T) {
			err := tt.pkg.Remove()
			assert.Equal(t, tt.err, errors.Cause(err))
			fakeHandler, ok := tt.pkg.cmdHandler.(*fake.CommandHandler)
			if !ok {
				t.Fatalf("Something's wrong with the fakeHandler!")
			}
			assert.Equal(t, tt.binary, fakeHandler.RemovedBinaries[0])
		})
	}
}

func TestUpgradeTo(t *testing.T) {
	tests := []struct {
		pkg        Package
		newVersion string
		err        error
	}{
		{
			pkg: Package{
				&config.Package{
					URL:     "test.com/no/bla",
					Version: "v0.0.1",
				},
				fake.NewCommandHandler("go: extracting test.com/no/bla v0.0.2\n", "v0.0.2", nil),
			},
			newVersion: "v0.0.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.pkg.URL, func(t *testing.T) {
			err := tt.pkg.UpgradeTo(tt.newVersion)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.newVersion, tt.pkg.Version)
		})
	}

}
