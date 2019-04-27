package subcmds

import (
	"errors"
	"fmt"
	"os"
	"strings"

	types "git.ramonruettimann.ml/ramon/packago/app/apis/packago"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
	"github.com/spf13/cobra"
)

// NewCommandInstall creates a new instance of the
// install command
func NewCommandInstall(cfg *types.Configuration) *cobra.Command {
	var installName string
	cmd := &cobra.Command{
		Use:   "install [repo-url]",
		Short: "install packages from given repo",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Old Config: %v\n", cfg)
			// parse package
			newPackage, err := parsePackage(args[0], installName)
			if err != nil {
				fmt.Printf("Error parsing package: %v\n", err)
				os.Exit(-1)
			}

			// append to config
			err = cfg.Packages.Add(newPackage)
			if err != nil {
				if err.Error() == "package already in list" {
					fmt.Printf("Package already installed!\n")
					os.Exit(0)
				}
				fmt.Printf("Error adding package to index: %v\n", err)
				os.Exit(-1)
			}

			// call go get on repo
			err = newPackage.Install()
			if err != nil {
				fmt.Printf("Error installing package: %v\n", err)
				os.Exit(-1)
			}
			// write out config
			fmt.Printf("New Config: %v\n", cfg)
		},
	}

	cmd.Flags().StringVar(&installName, "install-name", "", "compile the binary to this name")

	return cmd
}

// parsePackage uses the url to parse the necessary information
// to create a package struct. This includes the package url,
// installName (binary name) and the version
func parsePackage(packageURL, installName string) (*packages.Package, error) {
	pkg := &packages.Package{}

	pkg.URL = packageURL[:strings.LastIndex(packageURL, "@")]

	lastURL := packageURL[strings.LastIndex(packageURL, "/")+1:]
	urlSplit := strings.Split(lastURL, "@")
	if len(urlSplit) == 0 {
		return nil, errors.New("No Package name provided in URL")
	}

	if installName != "" {
		pkg.InstallName = installName
	} else {
		pkg.InstallName = urlSplit[0]
	}
	if len(urlSplit) == 1 {
		// no version provided
		pkg.Version = "latest"
	} else {
		pkg.Version = urlSplit[1]
	}

	return pkg, nil
}
