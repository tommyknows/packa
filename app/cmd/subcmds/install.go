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
func NewCommandInstall(pkgH *packages.PackageHandler) *cobra.Command {
	return &cobra.Command{
		Use:   "install [repo-url][@version]",
		Short: "install packages from given repo",
		Long: `If no argument / repo-url is given, install all packages that are currently
defined in the config file.

If an argument / repo-url is given, install this package with the specified version (see
normal go get URL form). If no version is given, use latest by default`,
		Run: func(cmd *cobra.Command, args []string) {
			err := install(pkgH, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}
}

func install(pkgH *packages.PackageHandler, args []string) error {
	if len(args) == 0 {
		err := pkgH.InstallAll()
		if err != nil {
			return errors.Wrapf(err, "error installing all packages")
		}
		err = config.SaveConfig()
		if err != nil {
			return errors.Wrapf(err, "error writing package to config")
		}
	}

	for _, arg := range args {
		pkg := packages.CreatePackage(arg)
		err := pkgH.Install(pkg)
		if err != nil {
			if err == packages.ErrPackageAlreadyInstalled {
				fmt.Printf("%v: %v", pkg.URL, err)
				return nil
			}
			return errors.Wrapf(err, "could not install package")
		}
		// write out config
		err = config.SavePackages(pkgH.ExportPackages())
		if err != nil {
			return errors.Wrapf(err, "error writing package to config")
		}
	}
	return nil
}
