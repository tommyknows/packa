package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"git.ramonruettimann.ml/ramon/packa/app/apis/defaults"
	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

var (
	// cfg is the Configuration initialised
	// to its defaults
	cfg = &Configuration{
		filename: defaults.ConfigFileFullPath(),
		Packages: []*Package{
			{
				URL:     "git.ramonruettimann.ml/ramon/packa",
				Version: "latest",
			},
		},
	}
	// Default provides default values for the configuration
	Default = Configuration{}
)

// Configuration is the root element for the config file
type Configuration struct {
	// Config is the configuration for the program
	Config AppConfig `yaml:"Config,omitempty"`
	// Packages contains the full list of packages in the state file
	Packages []*Package `yaml:"Packages,omitempty"`
	// filename is just a temporary info to write the config
	// back to the original source
	filename string
}

// Package contains info about a package that needs to be
// installed
type Package struct {
	// URL where to get the package from
	URL string `yaml:"URL"`
	// Which version should be installed (semver, go modules!)
	Version string `yaml:"Version"`
}

// AppConfig specifies configuration for the packa application
type AppConfig struct {
	// workingDir specifies where all the go get commands are executed
	WorkingDir string `yaml:"WorkingDir,omitempty"`
	// binaryDir specifies the binary location (defaults to $GOPATH/bin)
	BinaryDir string `yaml:"BinaryDir,omitempty"`
}

// WorkingDir gets the default working directory
func WorkingDir() string {
	if cfg.Config.WorkingDir != "" {
		return cfg.Config.WorkingDir
	}
	return defaults.WorkingDir()
}

// BinaryDir returns the directory in which to
// search for installed Go Binaries
func BinaryDir() string {
	if cfg.Config.BinaryDir != "" {
		return cfg.Config.BinaryDir
	}
	return defaults.BinaryDir()
}

// Load the config from cfgFile into cfg
func Load(cfgFile string) Configuration {
	// if cfgFile is not defined, get the default config file name
	if cfgFile == "" {
		cfgFile = defaults.ConfigFileFullPath()
		// create directory if not exists
		if _, err := os.Stat(path.Dir(cfgFile)); os.IsNotExist(err) {
			err := os.MkdirAll(path.Dir(cfgFile), 0777)
			if err != nil {
				klog.Fatalf("could not create directory for config file: %v", err)
			}
			klog.Infof("Created default working directory at %v", path.Dir(cfgFile))
		}
	}
	klog.V(3).Infof("Using config file from %v", cfgFile)

	//If a config file is found, read it in.
	contents, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			// returns the default config
			klog.Infof("returning default config as file does not yet exist")
			return *cfg
		}
		klog.Fatalf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(contents, cfg)
	if err != nil {
		klog.Fatalf("Cannot continue without a valid config file: %v", err)
	}
	cfg.filename = cfgFile
	klog.V(3).Infof("Parsed config file: %v\n", cfg)
	return *cfg
}

// SavePackages to file
func SavePackages(pkgs []*Package) error {
	cfg.Packages = pkgs
	return SaveConfig()
}

// SaveConfig to location
func SaveConfig() error {
	if cfg == nil {
		return errors.New("no config to save")
	}
	klog.V(3).Infof("Saving config to %v", cfg.filename)
	y, err := yaml.Marshal(cfg)
	if err != nil {
		klog.Fatalf("Error marshaling config to yaml: %v", err)
	}

	err = ioutil.WriteFile(cfg.filename, y, 0644)
	if err != nil {
		klog.Fatalf("error writing config file: %v", err)
	}
	klog.V(3).Infof("wrote config back to %v", cfg.filename)
	return nil
}
