package app

import (
	"flag"

	"git.ramonruettimann.ml/ramon/packa/app/cmd"
	"github.com/spf13/pflag"
)

// Run creates and executes new packa command
func Run() error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	_ = pflag.Set("logtostderr", "true")
	//We do not want these flags to show up in --help
	//These MarkHidden calls must be after the lines above
	_ = pflag.CommandLine.MarkHidden("version")
	_ = pflag.CommandLine.MarkHidden("log-flush-frequency")
	_ = pflag.CommandLine.MarkHidden("alsologtostderr")
	_ = pflag.CommandLine.MarkHidden("log-backtrace-at")
	_ = pflag.CommandLine.MarkHidden("log-dir")
	_ = pflag.CommandLine.MarkHidden("logtostderr")
	_ = pflag.CommandLine.MarkHidden("stderrthreshold")
	_ = pflag.CommandLine.MarkHidden("vmodule")

	cmd := cmd.NewPackaCommand()
	return cmd.Execute()
}
