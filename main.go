package main

import (
	"ead-tools/cmd"
	"fmt"
	"os"
)

func main() {
	fmt.Print("== ead-tools")
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error(), "\n")
		cmd.CloseLog()
		os.Exit(1)
	}
	cmd.CloseLog()
	cmd.ScanLog()
	os.Exit(0)
}
