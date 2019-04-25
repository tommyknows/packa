package main

import (
	"fmt"
	"os"

	"git.ramonruettimann.ml/ramon/packago/app"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
