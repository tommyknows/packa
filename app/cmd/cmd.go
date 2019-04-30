// Package cmd implements all cobra-command functions for packago
// This means initialising the main command and parsing the config file
package cmd

import (
	"fmt"
	"io"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"git.ramonruettimann.ml/ramon/packago/app/cmd/subcmds"
	"git.ramonruettimann.ml/ramon/packago/pkg/packages"
	"github.com/spf13/cobra"
)

// NewPackagoCommand returns a cobra command with default parameters
func NewPackagoCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var cfg config.Configuration
	var cfgFile string
	var pkgH packages.PackageHandler
	cmd := &cobra.Command{
		Version: version,
		Use:     "packago",
		Short:   "packago is a package manager for go",
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")
	cobra.OnInitialize(func() {
		cfg = config.Load(cfgFile)
		pkgHP, err := packages.NewPackageHandler(packages.Handle(cfg.Packages))
		pkgH = *pkgHP
		if err != nil {
			fmt.Printf("Error setting up CLI: %v\n", err)
			os.Exit(-1)
		}
	})
	cmd.AddCommand(subcmds.NewCommandInstall(&pkgH))
	cmd.AddCommand(subcmds.NewCommandUpgrade(&pkgH))
	cmd.AddCommand(subcmds.NewCommandRemove(&pkgH))

	return cmd
}
