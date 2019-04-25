// Package cmd implements all cobra-command functions for packago
// This means initialising the main command and parsing the config file
package cmd

import (
	"io"
	"os"

	types "git.ramonruettimann.ml/ramon/packago/app/apis/packago"
	"git.ramonruettimann.ml/ramon/packago/app/cmd/subcmds"
	"git.ramonruettimann.ml/ramon/packago/app/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"
)

var (
	cfgFile string
	config  types.Configuration
)

// NewPackagoCommand returns a cobra command with default parameters
func NewPackagoCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Version: version,
		Use:     "packago",
		Short:   "packago is a package manager for go",
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")
	cobra.OnInitialize(initConfig)
	cmd.AddCommand(subcmds.NewCommandInstall(&config))

	return cmd

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Read in default config values
	config = types.DefaultConfig

	if cfgFile != "" {
		// Use config file from the flag.
		klog.V(1).Infof("Setting config file to %v", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(constants.ConfigFileName)
		viper.AddConfigPath(constants.ConfigFileLocalDir)
	}

	//If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// TODO
			klog.Fatalf("NOT IMPLEMENTED: CREATING CONFIG FILE")
		default:
			klog.Fatalf("Unknown error occured while reading config: %v", err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		klog.Errorf("Cannot continue without a valid config file: %v", err)
		os.Exit(1)
	}
	klog.V(3).Infof("Parsed config file: %v\n", config)
}
