package main

import (
	"os"

	"git.ramonruettimann.ml/ramon/packa/app"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
