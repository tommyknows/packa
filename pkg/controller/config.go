package controller

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

type Configuration struct {
	// Settings for packa and its handlers
	Settings *Settings `json:"settings"`
	// all the handlers and their packages
	Packages map[string]*json.RawMessage `json:"packages,omitempty"`
	// for operations on the config file (save / close)
	file string
}

type Settings struct {
	// Settings for Handlers
	Handler map[string]*json.RawMessage `json:"handler,omitempty"`
}

// defaultConfig returns the default config that
// is save to use with the controller
func defaultConfig() *Configuration {
	return &Configuration{
		Packages: make(map[string]*json.RawMessage),
		Settings: &Settings{
			Handler: make(map[string]*json.RawMessage),
		},
	}
}

// Option for the controller initialisation.
// Directly supply the configuration. This assumes that all fields
// are at least initialised!
func Config(cfg *Configuration) Option {
	return func(ctl *Controller) error {
		ctl.configuration = cfg
		return nil
	}
}

// Option for the controller initialisation.
// ConfigFile reads in the configuration from the given file location
func ConfigFile(cfgFile string) Option {
	return func(ctl *Controller) error {
		f, err := os.OpenFile(cfgFile, os.O_RDONLY, os.ModeTemporary)
		if err != nil {
			return errors.Wrapf(err, "could not open config file")
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.Wrapf(err, "could not read config file")
		}
		ctl.configuration = defaultConfig()

		err = yaml.Unmarshal(data, ctl.configuration)
		if err != nil {
			return errors.Wrapf(err, "could not unmarshal")
		}

		ctl.configuration.file = cfgFile

		return err
	}
}

var errorFileNotSet = errors.New("no file has been set")

// save the config file to the file, if set. If no File
// should be set, save returns errorFileNotSet
func (cfg *Configuration) save() error {
	if cfg.file == "" {
		return errorFileNotSet
	}

	f, err := os.OpenFile(cfg.file, os.O_WRONLY, os.ModeTemporary)
	if err != nil {
		return errors.Wrapf(err, "could not open config file")
	}
	defer f.Close()

	enc, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrapf(err, "could not marshal config file")
	}

	err = f.Truncate(0)
	if err != nil {
		return errors.Wrapf(err, "could not clear config file to overwrite")
	}

	_, _ = f.Seek(0, 0)
	bw, err := f.Write(enc)
	if err != nil {
		return errors.Wrapf(err, "could not write config file")
	}
	if bw != len(enc) {
		return errors.Errorf("did not write correct config to config file, expected to write %v bytes, wrote %v bytes", len(enc), bw)
	}
	return nil
}
