package cmd

import (
	"fmt"
	"os"
)

func UseColor() bool {
	finfo, err := os.Stdout.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	return finfo.Mode()&os.ModeCharDevice == os.ModeCharDevice
}
