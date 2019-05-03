package app

import (
	"flag"
	"os"

	"git.ramonruettimann.ml/ramon/packa/app/cmd"
	"github.com/spf13/pflag"
)

// Run creates and executes new nine-access-guard command
func Run() error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Set("logtostderr", "true")
	//We do not want these flags to show up in --help
	//These MarkHidden calls must be after the lines above
	pflag.CommandLine.MarkHidden("version")
	pflag.CommandLine.MarkHidden("log-flush-frequency")
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log-backtrace-at")
	pflag.CommandLine.MarkHidden("log-dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("vmodule")

	cmd := cmd.NewPackagoCommand(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
