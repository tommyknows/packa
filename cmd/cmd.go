package cmd

import (
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tommyknows/packa/pkg/controller"
	"github.com/tommyknows/packa/pkg/defaults"
	"github.com/tommyknows/packa/pkg/handlers/brew"
	"github.com/tommyknows/packa/pkg/handlers/goget"
	"github.com/tommyknows/packa/pkg/output"
	"k8s.io/klog"
)

// If we did some kind of official binary release, we could
// set this version in the pipeline to the actual release
// version and / or commit id. but we don't, so indicate that
// people grabbed this from master
var version = "master"

const Name = "packa"

type PackageHandler interface {
	controller.PackageHandler
	Name() string
	Command() *cobra.Command
}

// NewPackaCommand returns the root command for packa
func NewPackaCommand() *cobra.Command {
	var cfgFile string
	var ctl *controller.Controller
	cmd := &cobra.Command{
		Version:      version,
		Use:          Name,
		Short:        "packa is a package manager",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")

	// if cfgFile is not defined, get the default config file name
	if cfgFile == "" {
		var err error
		cfgFile, err = createConfigFile()
		if err != nil {
			klog.Fatalf("could not create default config file location: %v", err)
		}
	}

	h := make(map[string]controller.PackageHandler)
	for _, handler := range []PackageHandler{goget.New(), brew.New()} {
		h[handler.Name()] = handler
	}

	ctl, err := controller.New(
		controller.ConfigFile(cfgFile),
		controller.RegisterHandlers(h),
	)
	if err != nil {
		klog.Fatalf("could not create controller: %v", err)
	}

	subcmds := []func(*controller.Controller) *cobra.Command{
		installCommand,
		upgradeCommand,
		removeCommand,
		listCommand,
	}

	for _, handler := range h {
		// this conversion is save as we cast it from a
		// PackageHandler when adding it to the map.
		c := handler.(PackageHandler).Command()
		for _, sub := range subcmds {
			c.AddCommand(sub(ctl))
		}
		cmd.AddCommand(c)
	}

	// these commands can be called without a handler
	// to do an operation on all of them.
	cmd.AddCommand(upgradeCommand(ctl))
	cmd.AddCommand(listCommand(ctl))

	return cmd
}

// createConfigFileLocation creates the config file
// directory and the file itself, if they should not
// exist already, and then returns the path to the file
func createConfigFile() (cfgFilePath string, err error) {
	// if cfgFile is not defined, get the default config file name
	cfgFile := defaults.ConfigFileFullPath()
	// create directory if not exists
	if _, err := os.Stat(path.Dir(cfgFile)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(cfgFile), 0777)
		if err != nil {
			return "", errors.Wrapf(err, "could not create default directory for config file")
		}
		klog.Infof("Created default working directory at %v", path.Dir(cfgFile))
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		f, err := os.Create(cfgFile)
		f.Close()
		if err != nil {
			return "", errors.Wrapf(err, "could not create empty config file")
		}
		klog.Infof("Created empty config file at %v", cfgFile)
	}
	return cfgFile, nil
}

// always try to close the controller, this will save the config
// file. if an error occurred on closing but we are already returning
// a non-nil error, we just log it.
// if no err would be returned normally, we overwrite it
// use it with:
//  defer func() {
//    err = close(ctl, err)
//  }()
func close(ctl io.Closer, inerr error) (outerr error) {
	outerr = inerr
	if err := ctl.Close(); err != nil {
		if inerr != nil {
			output.Warn(err.Error())
			return nil
		}
		outerr = err
	}
	return outerr
}
