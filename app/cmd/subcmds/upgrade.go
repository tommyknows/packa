package subcmds

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/app/apis/config"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	packages "git.ramonruettimann.ml/ramon/packa/pkg/packagehandler"
	"github.com/spf13/cobra"
)

// NewCommandUpgrade creates a new instance of the
// upgrade command
func NewCommandUpgrade(pkgH *packages.PackageHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [repo-url]",
		Short: "upgrades packages from given repo",
		Long: `Upgrade all packages that have "latest"
or "master" set as their versions. Packages pinned to a
specific version will not be touched. To update these,
call "install" with the given package and version you'd
like to upgrade to.`,
		Args: cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			err := upgrade(pkgH, args)
			if err != nil {
				output.Error(err.Error())
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
			output.Error("Could not write package state: %v", err)
		}
	}()

	return pkgH.UpgradeAll()
}
