package subcmds

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
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
			switch len(args) {
			case 0:
				(*cfg).Packages.InstallAll()
			case 1:
				pkg := packages.CreatePackage(args[0])
				err := (*cfg).Packages.Install(pkg)
				if err != nil {
					if err == packages.ErrPackageAlreadyInstalled {
						fmt.Printf("%v: %v", pkg.URL, err)
						os.Exit(0)
					}
					fmt.Printf("Could not install package: %v", err)
					os.Exit(-1)
				}

			}
			// write out config
			err := cfg.SaveConfig()
			if err != nil {
				fmt.Printf("error saving config: %v", err)
				os.Exit(-1)
			}
		},
	}

	cmd.Flags().StringVar(&installName, "install-name", "", "compile the binary to this name")

	return cmd
}
