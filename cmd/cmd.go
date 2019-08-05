package cmd

import (
	"os"
	"path"

	"git.ramonruettimann.ml/ramon/packa/cmd/subcmd"
	"git.ramonruettimann.ml/ramon/packa/pkg/controller"
	"git.ramonruettimann.ml/ramon/packa/pkg/defaults"
	"git.ramonruettimann.ml/ramon/packa/pkg/handlers/goget"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// NewPackaCommand creates the root packa command, adds all subcommands
// and initiates the controller
func NewPackaCommand() *cobra.Command {
	var cfgFile string
	var ctl *controller.Controller
	cmd := &cobra.Command{
		Version: version,
		Use:     "packa",
		Short:   "packa is a package manager",
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")

	// if cfgFile is not defined, get the default config file name
	if cfgFile == "" {
		var err error
		cfgFile, err = createConfigLocation()
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

func createConfigLocation() (string, error) {
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
