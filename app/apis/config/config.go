package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

var cfg *Configuration

const (
	// ConfigFileLocalDir is the user-local directory to search for
	// a config file
	ConfigFileLocalDir = "$HOME/.packago"
	// ConfigFileName is the default name for the config file
	ConfigFileName = "packago.yml"
)

// GetWorkingDir gets the default working directory
func GetWorkingDir() string {
	if cfg.Config.WorkingDir != "" {
		return cfg.Config.WorkingDir
	}
	return path.Join(os.Getenv("HOME"), ".packago")
}

// GetBinaryDir returns the default directory in which to
// search for installed Go Binaries
func GetBinaryDir() string {
	// if cfg ??
	// return cfg value if set
	return path.Join(os.Getenv("GOPATH"), "bin")
}

// Configuration is the root element for the config file
type Configuration struct {
	Config   appConfig  `yaml:"Config,omitempty"`
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
	// internal: InstalledVersion
	InstalledVersion string `yaml:"InstalledVersion,omitempty"`
}

type appConfig struct {
	// workingDir specifies where all the go get commands are executed
	WorkingDir string `yaml:"WorkingDir,omitempty"`
}

var (
	// Default provides default values for the configuration
	Default = Configuration{}
)

// Load the config from cfgFile into cfg
func Load(cfgFile string) Configuration {
	// Read in default config values
	cfg = &Default

	// if cfgFile is not defined, get the default config file name
	if cfgFile == "" {
		cfgFile = path.Join(ConfigFileLocalDir, ConfigFileName)
	}
	klog.V(3).Infof("Using config file from %v", cfgFile)

	//If a config file is found, read it in.
	contents, err := ioutil.ReadFile(cfgFile)
	if err != nil && !os.IsNotExist(err) {
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
	y, err := yaml.Marshal(cfg)
	if err != nil {
		klog.Fatalf("Error marshaling config to yaml: %v", err)
	}
	fmt.Printf("%s", y)
	err = ioutil.WriteFile(cfg.filename, y, 0644)
	if err != nil {
		klog.Fatalf("error writing config file: %v", err)
	}
	klog.V(3).Infof("wrote config back to %v", cfg.filename)
	return nil
}
