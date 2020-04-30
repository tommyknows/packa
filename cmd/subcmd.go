package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tommyknows/packa/pkg/controller"
	"github.com/tommyknows/packa/pkg/output"
)

func listCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list packages of handlers",
		Long: `if called with zero arguments, list will output
all packages of all handlers.
If called with one or more arguments, it will print the packages
of all specified handlers.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// always try to close the controller, this will save the config
			// file. if an error occurred on closing but we are already returning
			// a non-nil error, we just log it.
			// if no err would be returned normally, it is overwritten.
			defer func() {
				err = close(ctl, err)
			}()

			if cmd.Parent().Name() == Name {
				return ctl.PrintPackages(args...)
			}
			return ctl.PrintPackages(cmd.Parent().Name())
		},
	}
}

func installCommand(ctl *controller.Controller, handlers ...*cobra.Command) *cobra.Command {
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

			return ctl.Install(cmd.Parent().Name(), args...)
		},
	}
}

func removeCommand(ctl *controller.Controller) *cobra.Command {
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

			if len(args) == 0 {
				output.Warn("no package specified!")
				return nil
			}

			return ctl.Remove(cmd.Parent().Name(), args...)
		},
	}
}

func upgradeCommand(ctl *controller.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [package]",
		Short: "upgrade packages with the specified handler",
		Long: `upgrade a package with the specified handler. The package
name is handler-specific, check the documentation of the handler to get
the correct format. If no package name is given, upgrade all packages
that are in the index.
If no handler is set, upgrade all packages of all handlers.
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				err = close(ctl, err)
			}()

			// if the parent is not a handler, we want to upgrade all handlers
			if cmd.Parent().Name() == Name && len(args) == 0 {
				return ctl.UpgradeAll()
			}
			return ctl.Upgrade(cmd.Parent().Name(), args...)
		},
	}
}
