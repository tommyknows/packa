package subcmd

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/pkg/controller"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [handler] [package]",
		Short: "remove packages with the specified handler",
		Long: `remove a package with the specified handler. The package
name is handler-specific, check the documentation of the handler to get
the correct format.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			// exitOne temporary stores if we want to exit with return
			// code 1, as we still
			exitOne := false

			switch len(args) {
			case 2:
				err = ctl.Remove(args[0], args[1:]...)
			default:
				output.Warn("either no package or handler has been specified")
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
