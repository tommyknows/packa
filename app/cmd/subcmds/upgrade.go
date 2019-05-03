package subcmds

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	packages "git.ramonruettimann.ml/ramon/packa/pkg/packagehandler"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// NewCommandUpgrade creates a new instance of the
// upgrade command
func NewCommandUpgrade(pkgH *packages.PackageHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [repo-url]",
		Short: "upgrades packages from given repo",
		Long: `If no argument is given, all packages that have "latest"
or "master" set as their versions will be upgraded. Packages pinned to a
specific version will not be touched.

If a repo-url is given, update the given package to the specified version`,
		Args: cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			err := upgrade(pkgH, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}

	return cmd
}

func upgrade(pkgH *packages.PackageHandler, args []string) error {
	defer func() {
		err := config.SavePackages(pkgH.ExportPackages())
		if err != nil {
			klog.Fatalf("Could not write package state: %v", err)
		}
	}()

	return pkgH.UpgradeAll()
}
