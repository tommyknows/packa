package subcmds

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCommandRemove creates a new instance of the
// upgrade command
func NewCommandRemove(pkgH *packages.PackageHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [repo-url]",
		Short: "remove packages from given repo",
		Long:  `Remove the given package / repo-url from the index and remove its binary`,
		Run: func(cmd *cobra.Command, args []string) {
			err := remove(pkgH, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}

	return cmd
}

func remove(pkgH *packages.PackageHandler, args []string) error {
	if len(args) == 0 {
		return errors.New("nothing to delete, please specify at least one package to delete")
	}

	for _, arg := range args {
		p := packages.CreatePackage(arg)
		pkg := pkgH.GetPackage(p.URL)
		if pkg.Package == nil {
			return errors.New("Package " + arg + " is not installed")
		}
		err := pkgH.Remove(pkg)
		if err != nil {
			return fmt.Errorf("could not remove package:\n%v", err)
		}
		// write out config
		err = config.SaveConfig()
		if err != nil {
			return errors.Wrapf(err, "could not write config")
		}
	}
	return nil
}
