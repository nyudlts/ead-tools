package cmd

import (
	"github.com/spf13/cobra"
	"regexp"
)

var (
	inputDir string
)

var xmlPtn = regexp.MustCompile("xml$")

var rootCmd = &cobra.Command{}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
