package subcmds

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	packages "git.ramonruettimann.ml/ramon/packa/pkg/packagehandler"
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
				output.Error(err.Error())
				os.Exit(-1)
			}
		},
	}

	return cmd
}

func remove(pkgH *packages.PackageHandler, args []string) error {
	defer func() {
		err := config.SavePackages(pkgH.ExportPackages())
		if err != nil {
			output.Error("Could not write package state: %v", err)
		}
	}()

	if len(args) == 0 {
		return errors.New("nothing to delete, please specify at least one package to delete")
	}

	err := pkgH.Remove(pkgH.GetPackages(args...)...)
	return errors.Wrapf(err, "could not remove package(s)")
}
