package subcmd

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/pkg/controller"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	"github.com/spf13/cobra"
)

func NewInstallCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "install [handler] [package]",
		Short: "install packages with the specified handler",
		Long: `install a package with the specified handler. The package
name is handler-specific, check the documentation of the handler to get
the correct format. If no package name is given, install all packages
that are defined in the index.
Will also add the package to the index if it does not exist yet.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			// exitOne temporary stores if we want to exit with return
			// code 1, as we still want to close the controller and save
			// the config
			exitOne := false

			switch len(args) {
			case 0:
				output.Error("need to specifiy at least the handler name")
				os.Exit(-1)
			case 1:
				err = ctl.Install(args[0])
			default:
				err = ctl.Install(args[0], args[1:]...)
			}
			if err != nil {
				output.Error(err.Error())
				exitOne = true
			}
			err = ctl.Close()
			if err != nil {
				output.Error(err.Error())
				exitOne = true
			}
			if exitOne {
				os.Exit(-1)
			}
		},
	}
}
