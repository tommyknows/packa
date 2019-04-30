package subcmds

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCommandInstall creates a new instance of the
// install command
func NewCommandInstall(cfg *config.Configuration) *cobra.Command {
	var installName string
	cmd := &cobra.Command{
		Use:   "install [repo-url]",
		Short: "install packages from given repo",
		Run: func(cmd *cobra.Command, args []string) {
			err := install(cfg, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}

	cmd.Flags().StringVar(&installName, "install-name", "", "compile the binary to this name")

	return cmd
}

func install(cfg *config.Configuration, args []string) error {
	if len(args) == 0 {
		err := (*cfg).Packages.InstallAll()
		if err != nil {
			return errors.Wrapf(err, "error installing all packages")
		}
		err = cfg.SaveConfig()
		if err != nil {
			return errors.Wrapf(err, "error writing package to config")
		}
	}

	for _, arg := range args {
		pkg := packages.CreatePackage(arg)
		err := (*cfg).Packages.Install(pkg)
		if err != nil {
			if err == packages.ErrPackageAlreadyInstalled {
				fmt.Printf("%v: %v", pkg.URL, err)
				return nil
			}
			return errors.Wrapf(err, "could not install package")
		}
		// write out config
		err = cfg.SaveConfig()
		if err != nil {
			return errors.Wrapf(err, "error writing package to config")
		}

	}
	return nil
}
