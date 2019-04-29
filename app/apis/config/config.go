package config

import (
	"io/ioutil"
	"os"
	"path"

	"git.ramonruettimann.ml/ramon/packago/app/constants"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

// Configuration is the root element for the config file
type Configuration struct {
	Packages packages.Packages `mapstructure:"packages"`
	filename string
}

var (
	// Default provides default values for the configuration
	Default = Configuration{}
)

// Load the config from cfgFile into cfg
func Load(cfgFile string) Configuration {
	// Read in default config values
	cfg := &Default

	if cfgFile == "" {
		cfgFile = path.Join(constants.GetDefaultWorkingDir(), constants.ConfigFileName)
	}
	klog.V(3).Infof("Using config file from %v", cfgFile)

	//If a config file is found, read it in.
	//if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
	//}
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

// SaveConfig to location
func (cfg *Configuration) SaveConfig() error {
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
