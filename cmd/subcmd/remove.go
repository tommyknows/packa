package subcmd

import (
	"github.com/spf13/cobra"
	"github.com/tommyknows/packa/pkg/controller"
	"github.com/tommyknows/packa/pkg/output"
)

func NewRemoveCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <handler> <package>",
		Short: "remove packages with the specified handler",
		Long: `remove a package with the specified handler. The package
name is handler-specific, check the documentation of the handler to get
the correct format.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				err = close(ctl, err)
			}()
			switch len(args) {
			case 2:
				return ctl.Remove(args[0], args[1:]...)
			default:
				output.Warn("either no package or handler has been specified")
				return nil
			}
		},
	}
}
