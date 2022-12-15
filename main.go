package main

import (
	"ead-tools/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error(), "\n")
		os.Exit(1)
	}
	os.Exit(0)
}
