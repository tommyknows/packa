package subcmd

import (
	"errors"

	"github.com/tommyknows/packa/pkg/controller"
	"github.com/spf13/cobra"
)

func NewInstallCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "install <handler> [package]",
		Short: "install packages with the specified handler",
		Long: `install a package with the specified handler. The package
name is handler-specific, check the documentation of the handler to get
the correct format. If no package name is given, install all packages
that are defined in the index.
Will also add the package to the index if it does not exist yet.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				err = close(ctl, err)
			}()

			switch len(args) {
			case 0:
				return errors.New("need to specify at least the handler name")
			case 1:
				return ctl.Install(args[0])
			default:
				return ctl.Install(args[0], args[1:]...)
			}
		},
	}
}
