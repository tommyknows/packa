package packages

import (
	"fmt"
	"testing"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"git.ramonruettimann.ml/ramon/packa/test/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	packages := []*config.Package{
		{
			URL:     "test.com/test/test",
			Version: "v0.0.1",
		},
		{
			URL:     "test.com/another/secondtest",
			Version: latest,
		},
	}
	pkgH, err := NewPackageHandler(
		fake.NewCommandHandler("", nil),
		Handle(packages),
	)
	assert.Nil(t, err)
	assert.Equal(t, packages, pkgH.ExportPackages())
}

func TestGetPackages(t *testing.T) {
	packages := []*config.Package{
		{
			URL:     "test.com/test/test",
			Version: "v0.0.1",
		},
	}
	pkgH, err := NewPackageHandler(nil, Handle(packages))
	assert.Nil(t, err)

	pkgs := pkgH.GetPackages("test.com/test/test@v0.0.1")
	assert.Equal(t, pkgs, []Package{Package{packages[0], nil}})

	p := Package{
		&config.Package{
			URL:     "new.test/new",
			Version: "latest",
		},
		nil,
	}
	pkgs = pkgH.GetPackages(p.URL)

	assert.Equal(t, pkgs, []Package{p})
}

func TestInstallPackages(t *testing.T) {
	failCmdInit := func() CommandHandler {
		return fake.NewCommandHandler("error 123\n", fmt.Errorf("exit code 1"))
	}
	successCmdInit := func() CommandHandler {
		return fake.NewCommandHandler("go: extracting test.com/test/test v0.0.1\n", nil)
	}
	pkg := []*config.Package{
		{

			URL:     "test.com/test/test",
			Version: "v0.0.1",
		},
	}

	tests := []struct {
		name string
		cmdH CommandHandler
		// packages that will be added to the packagehandler
		alreadyInstalled []*config.Package
		// packages that will be used to call the install command
		toInstall []*config.Package
		// optional
		err func(CommandHandler) error
	}{
		{
			name:             "successful update",
			cmdH:             successCmdInit(),
			alreadyInstalled: pkg,
			toInstall:        pkg,
			err:              nil,
		},
		{
			name:             "successful install all",
			cmdH:             successCmdInit(),
			alreadyInstalled: pkg,
			toInstall:        []*config.Package{},
			err:              nil,
		},
		{
			name:             "successful update",
			cmdH:             successCmdInit(),
			alreadyInstalled: []*config.Package{},
			toInstall:        pkg,
			err:              nil,
		},
		{
			name:             "error installation",
			cmdH:             failCmdInit(),
			alreadyInstalled: pkg,
			toInstall:        pkg,
			err: func(cmdH CommandHandler) error {
				collErr := make(InstallError)
				collErr[Package{pkg[0], cmdH}] = errors.Wrapf(fmt.Errorf("exit code 1"), "could not install package %v", pkg[0].URL)
				return collErr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgH, err := NewPackageHandler(tt.cmdH, Handle(tt.alreadyInstalled))
			assert.Nil(t, err)

			err = pkgH.Install(
				func() []Package {
					pkgs := []Package{}
					for _, p := range tt.toInstall {
						pkgs = append(pkgs, Package{p, pkgH.cmdHandler})
					}
					return pkgs
				}()...,
			)

			if tt.err != nil {
				assert.Equal(t, tt.err(tt.cmdH).Error(), err.Error())
			} else {
				if err != nil {
					t.Errorf("Did not expect error, but got %v", err)
				}
			}

			fakeH, ok := tt.cmdH.(*fake.CommandHandler)
			if !ok {
				t.Fatalf("Something's wrong with the fakeHandler!")
			}

			afterInstall := tt.toInstall
			if len(afterInstall) == 0 {
				afterInstall = tt.alreadyInstalled
			}

			assert.Equal(t,
				func() []string {
					pkgs := []string{}
					for _, p := range afterInstall {
						pkgs = append(pkgs, p.URL)
					}
					return pkgs
				}(),
				fakeH.InstalledPackages,
			)
		})
	}
}

func TestRemove(t *testing.T) {
	failCmdInit := func() CommandHandler {
		return fake.NewCommandHandler("no such file or directory\n", fmt.Errorf("exit code 1"))
	}
	successCmdInit := func() CommandHandler {
		return fake.NewCommandHandler("go: extracting test.com/test/test v0.0.1\n", nil)
	}

	pkg := []*config.Package{
		{

			URL:     "test.com/test/test",
			Version: "v0.0.1",
		},
	}

	tests := []struct {
		name             string
		cmdH             CommandHandler
		alreadyInstalled []*config.Package
		toRemove         []*config.Package
		removedBinaries  []string
		err              func(CommandHandler) error
	}{
		{
			name:             "not installed",
			cmdH:             successCmdInit(),
			alreadyInstalled: []*config.Package{},
			toRemove:         pkg,
			removedBinaries:  []string{},
			err: func(cmdH CommandHandler) error {
				collErr := make(InstallError)
				collErr[Package{pkg[0], cmdH}] = errors.Errorf("package test.com/test/test not installed")
				return collErr
			},
		},
		{
			name:             "normal removal",
			cmdH:             successCmdInit(),
			alreadyInstalled: pkg,
			toRemove:         pkg,
			removedBinaries:  []string{"test"},
			err:              nil,
		},
		{
			name:             "unsuccessful removal",
			cmdH:             failCmdInit(),
			alreadyInstalled: pkg,
			toRemove:         pkg,
			removedBinaries:  []string{"test"},
			err: func(cmdH CommandHandler) error {
				collErr := make(InstallError)
				collErr[Package{pkg[0], cmdH}] = errors.Errorf("error removing binary, not removing package %v from state file: could not remove package: exit code 1", pkg[0].URL)
				return collErr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgH, err := NewPackageHandler(tt.cmdH, Handle(tt.alreadyInstalled))
			assert.Nil(t, err)

			err = pkgH.Remove(
				func() []Package {
					pkgs := []Package{}
					for _, p := range tt.toRemove {
						pkgs = append(pkgs, Package{p, pkgH.cmdHandler})
					}
					return pkgs
				}()...,
			)

			if tt.err != nil {
				assert.Equal(t, tt.err(tt.cmdH).Error(), err.Error())
			} else {
				if err != nil {
					t.Errorf("Did not expect error, but got %v", err)
				}
			}

			fakeH, ok := tt.cmdH.(*fake.CommandHandler)
			if !ok {
				t.Fatalf("Something's wrong with the fakeHandler!")
			}

			assert.Equal(t, tt.removedBinaries, fakeH.RemovedBinaries)
		})
	}
}

func TestUpgradeAll(t *testing.T) {
	failCmdInit := func() CommandHandler {
		return fake.NewCommandHandler("no such file or directory\n", fmt.Errorf("exit code 1"))
	}
	successCmdInit := func() CommandHandler {
		return fake.NewCommandHandler("go: extracting test.com/test/test v0.0.1\n", nil)
	}

	pkg0 := &config.Package{
		URL:     "test.com/test/test",
		Version: "latest",
	}

	tests := []struct {
		name             string
		cmdH             CommandHandler
		alreadyInstalled []*config.Package
		update           []*config.Package
		err              func(CommandHandler) error
	}{
		{
			name: "sucessful upgrade all",
			cmdH: successCmdInit(),
			alreadyInstalled: []*config.Package{
				pkg0,
				{
					URL:     "github.com/test/bla",
					Version: "v0.0.1",
				},
			},
			update: []*config.Package{
				pkg0,
			},
			err: nil,
		},
		{
			name:             "error upgrade all",
			cmdH:             failCmdInit(),
			alreadyInstalled: []*config.Package{pkg0},
			update:           []*config.Package{pkg0},
			err: func(cmdH CommandHandler) error {
				collErr := make(InstallError)
				collErr[Package{pkg0, cmdH}] = errors.Errorf("package %v not upgraded: could not install package %v: exit code 1", pkg0.URL, pkg0.URL)
				return collErr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgH, err := NewPackageHandler(tt.cmdH, Handle(tt.alreadyInstalled))
			assert.Nil(t, err)

			err = pkgH.UpgradeAll()
			if tt.err != nil {
				assert.Equal(t, tt.err(tt.cmdH).Error(), err.Error())
			} else {
				if err != nil {
					t.Errorf("Did not expect error, but got %v", err)
				}
			}

			fakeH, ok := tt.cmdH.(*fake.CommandHandler)
			if !ok {
				t.Fatalf("Something's wrong with the fakeHandler!")
			}

			assert.Equal(t, func() []string {
				updated := []string{}
				for _, p := range tt.update {
					updated = append(updated, p.URL)
				}
				return updated
			}(), fakeH.InstalledPackages)
		})
	}
}
