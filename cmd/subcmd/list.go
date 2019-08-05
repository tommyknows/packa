package subcmd

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/pkg/controller"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	"github.com/spf13/cobra"
)

func NewListCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "list [handler]",
		Short: "list handler or packages of the defined handler",
		Long: `if called with zero arguments, list will output
a list of all available handlers.
else, it will print the packages of the specified handler `,
		Run: func(cmd *cobra.Command, args []string) {
			//switch len(args) {
			//case 0, 1:
			err := ctl.PrintPackages(args[0:]...)
			if err != nil {
				output.Error("could not print packages: %v", err)
			}
			//default:
			//output.Error("too many arguments specified (need one at max)")
			//}
			err = ctl.Close()
			if err != nil {
				output.Error(err.Error())
				os.Exit(-1)
			}
		},
	}
}
