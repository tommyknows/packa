package main

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/app"
	"git.ramonruettimann.ml/ramon/packa/pkg/output"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	if err := app.Run(); err != nil {
		output.Error("error: %v\n", err)
		os.Exit(1)
	}
}
