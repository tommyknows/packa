package controller

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/tommyknows/packa/test/fake"
)

func testConfig() *Configuration {
	return &Configuration{
		Settings: &Settings{
			Handler: map[string]*json.RawMessage{
				"fake": fake.DefaultSettingsRaw,
			},
		},
		Packages: map[string]*json.RawMessage{
			"fake": fake.EmptyPackages,
		},
	}
}

func TestNewAndInit(t *testing.T) {
	is := is.New(t)

	ctl, err := New(
		Config(testConfig()),
		RegisterHandlers(
			map[string]PackageHandler{
				"fake": &fake.Handler{},
			},
		),
	)

	is.NoErr(err)
	is.True(ctl.handlers["fake"].initialised == false)
	is.NoErr(ctl.initialiseHandler("fake"))
	fh, ok := ctl.handlers["fake"].PackageHandler.(*fake.Handler)
	is.True(ok)
	is.Equal(fake.DefaultSettings.WorkingDir, fh.Config.WorkingDir)
	is.True(ctl.handlers["fake"].initialised)
}

func TestInstall(t *testing.T) {
	is := is.New(t)

	fH := &fake.Handler{}
	ctl := &Controller{
		configuration: testConfig(),
		handlers: map[string]*handler{
			"fake": {
				fH, false,
			},
		},
	}

	err := ctl.Install("fake", "testpackage")
	is.NoErr(err)
	is.Equal(1, len(fH.Packages))
	is.Equal("testpackage", fH.Packages[0].Name)

	err = ctl.Install("nonexistenthandler", "test")
	is.Equal("handler \"nonexistenthandler\" does not exist or has not been registered", err.Error())
}

func TestRemove(t *testing.T) {
	is := is.New(t)

	cfg := testConfig()
	cfg.Packages["fake"] = fake.DefaultPackagesRaw

	fH := &fake.Handler{}
	ctl := &Controller{
		configuration: cfg,
		handlers: map[string]*handler{
			"fake": {
				fH, false,
			},
		},
	}

	err := ctl.Remove("fake", "fakePackage1")
	is.NoErr(err)
	is.Equal(1, len(fH.Packages))
	is.Equal("fakePackage2", fH.Packages[0].Name)
}

func TestUpgrade(t *testing.T) {
	is := is.New(t)
	cfg := testConfig()
	cfg.Packages["fake"] = fake.DefaultPackagesRaw

	fH := &fake.Handler{}
	ctl := &Controller{
		configuration: cfg,
		handlers: map[string]*handler{
			"fake": {
				fH, false,
			},
		},
	}

	err := ctl.Upgrade("fake", "fakePackage1")
	is.NoErr(err)
	is.Equal(2, len(fH.Packages))
	is.Equal("fakePackage1+", fH.Packages[0].Name)
	is.Equal("fakePackage2", fH.Packages[1].Name)
}

func TestUpgradeAll(t *testing.T) {
	is := is.New(t)

	cfg := testConfig()
	cfg.Packages["fake"] = fake.DefaultPackagesRaw
	cfg.Packages["fake2"] = fake.DefaultPackagesRaw

	fH1 := &fake.Handler{}
	fH2 := &fake.Handler{}
	ctl := &Controller{
		configuration: cfg,
		handlers: map[string]*handler{
			"fake": {
				fH1, false,
			},
			"fake2": {
				fH2, false,
			},
		},
	}

	err := ctl.UpgradeAll()
	is.NoErr(err)
	is.Equal(2, len(fH1.Packages))
	is.Equal("fakePackage1+", fH1.Packages[0].Name)
	is.Equal("fakePackage2+", fH1.Packages[1].Name)
	is.Equal(2, len(fH2.Packages))
	is.Equal("fakePackage1+", fH2.Packages[0].Name)
	is.Equal("fakePackage2+", fH2.Packages[1].Name)
}
