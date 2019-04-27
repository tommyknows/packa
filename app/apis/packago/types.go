package packago

import "git.ramonruettimann.ml/ramon/packago/pkg/packages"

// Configuration is the root element for the config file
type Configuration struct {
	Config   *Config           `mapstructure:"config"`
	Packages packages.Packages `mapstructure:"packages"`
}

// Config contains config options for packago
type Config struct {
	AutoUpdate bool `mapstructure:"autoUpdate"`
	// TODO: config options??
	// autoremove source code when removing packages
	AutoRemove bool `mapstructure:"autoRemove"`
}

var (
	// DefaultConfig provides default values for the configuration
	DefaultConfig = Configuration{
		Config: &Config{
			AutoUpdate: false,
		},
	}
)
