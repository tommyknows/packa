package subcmd

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/pkg/controller"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	"github.com/spf13/cobra"
)

func NewUpgradeCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [handler] [package]",
		Short: "upgrade packages with the specified handler",
		Long: `upgrade a package with the specified handler. The package
name is handler-specific, check the documentation of the handler to get
the correct format. If no package name is given, upgrade all packages
that are in the index.
`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			switch len(args) {
			case 0:
				err = ctl.UpgradeAll()
			case 1:
				err = ctl.Upgrade(args[0])
			default:
				err = ctl.Upgrade(args[0], args[1:]...)
			}
			if err != nil {
				output.Error(err.Error())
				os.Exit(-1)
			}
			err = ctl.Close()
			if err != nil {
				output.Error(err.Error())
				os.Exit(-1)
			}
		},
	}
}
