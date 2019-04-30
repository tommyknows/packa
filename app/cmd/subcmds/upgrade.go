package subcmds

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCommandUpgrade creates a new instance of the
// upgrade command
func NewCommandUpgrade(cfg *config.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [repo-url]",
		Short: "upgrades packages from given repo",
		Long: `If no argument is given, all packages that have "latest"
or "master" set as their versions will be upgraded. Packages pinned to a
specific version will not be touched.

If a repo-url is given, update the given package to the specified version`,

		Run: func(cmd *cobra.Command, args []string) {
			err := upgrade(cfg, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}

	return cmd
}

func upgrade(cfg *config.Configuration, args []string) error {
	if len(args) == 0 {
		(*cfg).Packages.UpgradeAll()
	}

	for _, arg := range args {
		p := packages.CreatePackage(arg)
		pkg := (*cfg).Packages.GetPackage(p.URL)
		if pkg == nil {
			return errors.New("Package " + arg + " is not installed")
		}
		// extract version and set the package to that,
		// so that install installs this version
		if pkg.Version == p.Version {
			return errors.New("Not upgrading package " + p.URL + "@" + p.Version + " as same version already installed")
		}
		pkg.Version = p.Version
		err := pkg.Install()
		if err != nil {
			if err == packages.ErrPackageAlreadyInstalled {
				fmt.Printf("%v: %v", pkg.URL, err)
				return nil
			}
			return fmt.Errorf("could not upgrade package:\n%v", err)
		}
	}
	// write out config
	return cfg.SaveConfig()
}
