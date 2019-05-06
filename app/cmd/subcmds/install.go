package subcmds

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	packages "git.ramonruettimann.ml/ramon/packa/pkg/packagehandler"
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
				output.Error(err.Error())
				os.Exit(-1)
			}
		},
	}
}

func install(pkgH *packages.PackageHandler, args []string) error {
	defer func() {
		err := config.SavePackages(pkgH.ExportPackages())
		if err != nil {
			output.Error("Could not write package state: %v", err)
		}
	}()

	err := pkgH.Install(pkgH.GetPackages(args...)...)
	return errors.Wrapf(err, "📦 could not install package(s)")
}
