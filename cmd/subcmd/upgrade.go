package subcmd

import (
	"git.ramonruettimann.ml/ramon/packa/pkg/controller"
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				err = close(ctl, err)
			}()

			switch len(args) {
			case 0:
				return ctl.UpgradeAll()
			case 1:
				return ctl.Upgrade(args[0])
			default:
				return ctl.Upgrade(args[0], args[1:]...)
			}
		},
	}
}
