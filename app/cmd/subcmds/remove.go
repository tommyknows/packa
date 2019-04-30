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
func NewCommandRemove(cfg *config.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [repo-url]",
		Short: "remove packages from given repo",
		Long:  `Remove the given package / repo-url from the index`,
		Run: func(cmd *cobra.Command, args []string) {
			err := remove(cfg, args)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}

	return cmd
}

func remove(cfg *config.Configuration, args []string) error {
	if len(args) == 0 {
		return errors.New("nothing to delete, please specify at least one package to delete")
	}

	for _, arg := range args {
		p := packages.CreatePackage(arg)
		pkg := (*cfg).Packages.GetPackage(p.URL)
		if pkg == nil {
			return errors.New("Package " + arg + " is not installed")
		}
		err := (*cfg).Packages.Remove(pkg)
		if err != nil {
			if err == packages.ErrPackageAlreadyInstalled {
				fmt.Printf("%v: %v", pkg.URL, err)
				return nil
			}
			return fmt.Errorf("could not upgrade package:\n%v", err)
		}
		// write out config
		err = cfg.SaveConfig()
		if err != nil {
			return errors.Wrapf(err, "could not write config")
		}
	}
	return nil
}
