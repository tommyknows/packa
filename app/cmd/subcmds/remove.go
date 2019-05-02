package subcmds

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	packages "git.ramonruettimann.ml/ramon/packago/pkg/packagehandler"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog"
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
	defer func() {
		err := config.SavePackages(pkgH.ExportPackages())
		if err != nil {
			klog.Fatalf("Could not write package state: %v", err)
		}
	}()

	if len(args) == 0 {
		return errors.New("nothing to delete, please specify at least one package to delete")
	}

	pkgs := pkgH.GetPackages(args...)
	err := pkgH.Remove(pkgs...)
	return errors.Wrapf(err, "could not remove some packages")
}
