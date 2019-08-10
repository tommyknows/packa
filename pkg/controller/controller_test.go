package controller

import (
	"encoding/json"
	"testing"

	"github.com/tommyknows/packa/test/fake"
	"github.com/stretchr/testify/assert"
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
	ctl, err := New(
		Config(testConfig()),
		RegisterHandlers(
			map[string]PackageHandler{
				"fake": &fake.Handler{},
			},
		),
	)

	assert.NoError(t, err)
	assert.False(t, ctl.handlers["fake"].initialised)
	assert.NoError(t, ctl.initialiseHandler("fake"))
	fh, ok := ctl.handlers["fake"].PackageHandler.(*fake.Handler)
	assert.True(t, ok)
	assert.Equal(t, fake.DefaultSettings.WorkingDir, fh.Config.WorkingDir)
	assert.True(t, ctl.handlers["fake"].initialised)
}

func TestInstall(t *testing.T) {
	fH := &fake.Handler{}
	ctl := &Controller{
		configuration: testConfig(),
		handlers: map[string]*handler{
			"fake": &handler{
				fH, false,
			},
		},
	}

	err := ctl.Install("fake", "testpackage")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fH.Packages))
	assert.Equal(t, "testpackage", fH.Packages[0].Name)

	err = ctl.Install("nonexistenthandler", "test")
	assert.Equal(t, "handler \"nonexistenthandler\" does not exist or has not been registered", err.Error())
}

func TestRemove(t *testing.T) {
	cfg := testConfig()
	cfg.Packages["fake"] = fake.DefaultPackagesRaw

	fH := &fake.Handler{}
	ctl := &Controller{
		configuration: cfg,
		handlers: map[string]*handler{
			"fake": &handler{
				fH, false,
			},
		},
	}

	err := ctl.Remove("fake", "fakePackage1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fH.Packages))
	assert.Equal(t, "fakePackage2", fH.Packages[0].Name)
}

func TestUpgrade(t *testing.T) {
	cfg := testConfig()
	cfg.Packages["fake"] = fake.DefaultPackagesRaw

	fH := &fake.Handler{}
	ctl := &Controller{
		configuration: cfg,
		handlers: map[string]*handler{
			"fake": &handler{
				fH, false,
			},
		},
	}

	err := ctl.Upgrade("fake", "fakePackage1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(fH.Packages))
	assert.Equal(t, "fakePackage1+", fH.Packages[0].Name)
	assert.Equal(t, "fakePackage2", fH.Packages[1].Name)
}

func TestUpgradeAll(t *testing.T) {
	cfg := testConfig()
	cfg.Packages["fake"] = fake.DefaultPackagesRaw
	cfg.Packages["fake2"] = fake.DefaultPackagesRaw

	fH1 := &fake.Handler{}
	fH2 := &fake.Handler{}
	ctl := &Controller{
		configuration: cfg,
		handlers: map[string]*handler{
			"fake": &handler{
				fH1, false,
			},
			"fake2": &handler{
				fH2, false,
			},
		},
	}

	err := ctl.UpgradeAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(fH1.Packages))
	assert.Equal(t, "fakePackage1+", fH1.Packages[0].Name)
	assert.Equal(t, "fakePackage2+", fH1.Packages[1].Name)
	assert.Equal(t, 2, len(fH2.Packages))
	assert.Equal(t, "fakePackage1+", fH2.Packages[0].Name)
	assert.Equal(t, "fakePackage2+", fH2.Packages[1].Name)
}
