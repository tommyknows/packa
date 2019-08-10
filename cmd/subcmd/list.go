package subcmd

import (
	"github.com/tommyknows/packa/pkg/controller"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewListCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "list [handler]",
		Short: "list handler or packages of the defined handler",
		Long: `if called with zero arguments, list will output
a list of all available handlers.
else, it will print the packages of the specified handler `,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// always try to close the controller, this will save the config
			// file. if an error occured on closing but we are already returning
			// a non-nil error, we just log it.
			// if no err would be returned normally, we overwrite it
			defer func() {
				err = close(ctl, err)
			}()
			return errors.Wrapf(ctl.PrintPackages(args[0:]...), "could not print packages")
		},
	}
}
