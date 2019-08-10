package main

import (
	"flag"
	"os"

	"github.com/tommyknows/packa/cmd"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
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
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
