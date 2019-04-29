// Package cmd implements all cobra-command functions for packago
// This means initialising the main command and parsing the config file
package cmd

import (
	"io"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"git.ramonruettimann.ml/ramon/packago/app/cmd/subcmds"
	"github.com/spf13/cobra"
)

// NewPackagoCommand returns a cobra command with default parameters
func NewPackagoCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var cfg config.Configuration
	var cfgFile string
	cmd := &cobra.Command{
		Version: version,
		Use:     "packago",
		Short:   "packago is a package manager for go",
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")
	cobra.OnInitialize(func() {
		cfg = config.Load(cfgFile)
	})
	cmd.AddCommand(subcmds.NewCommandInstall(&cfg))

	return cmd
}
