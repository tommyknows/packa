package cmd

import (
	"os"
	"path"

	"github.com/tommyknows/packa/cmd/subcmd"
	"github.com/tommyknows/packa/pkg/controller"
	"github.com/tommyknows/packa/pkg/defaults"
	"github.com/tommyknows/packa/pkg/handlers/goget"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// NewPackaCommand returns the root command for packa
func NewPackaCommand() *cobra.Command {
	var cfgFile string
	var ctl *controller.Controller
	cmd := &cobra.Command{
		Version:      version,
		Use:          "packa",
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

	ctl, err := controller.New(
		controller.ConfigFile(cfgFile),
		controller.RegisterHandlers(map[string]controller.PackageHandler{
			"go": goget.New(),
		}),
	)
	if err != nil {
		klog.Fatalf("could not create controller: %v", err)
	}

	cmd.AddCommand(subcmd.NewInstallCommand(ctl))
	cmd.AddCommand(subcmd.NewUpgradeCommand(ctl))
	cmd.AddCommand(subcmd.NewRemoveCommand(ctl))
	cmd.AddCommand(subcmd.NewListCommand(ctl))

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