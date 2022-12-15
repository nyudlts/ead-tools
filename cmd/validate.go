package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputDir, "input-dir", "i", ".", "the directory to scan")
	rootCmd.AddCommand(validateEadCmd)
}

var ignoreDirFilter = []string{".git", ".idea", ".circleci"}

func ignoreDir(name string) bool {
	for _, s := range ignoreDirFilter {
		if s == name {
			return true
		}
	}
	return false
}

var validateEadCmd = &cobra.Command{
	Use: "validate-ead",
	Run: func(cmd *cobra.Command, args []string) {

		if err := directoryExists(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(2)
		}

		log.Printf("[INFO] validating ead files in %s", inputDir)
		if err := validateEAD(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(3)
		}
	},
}

func directoryExists() error {
	if info, err := os.Stat(inputDir); err == nil {
		//check if it is a directory
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory")
		}

	} else {
		return err
	}
	return nil
}

func validateEAD() error {
	rootObjects, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, rootObject := range rootObjects {
		if !ignoreDir(rootObject.Name()) {
			if rootObject.IsDir() {
				if err = walkDirectory(filepath.Join(inputDir, rootObject.Name())); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func walkDirectory(directory string) error {
	log.Printf("[INFO] validating files in %s", directory)
	if err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			log.Printf("[INFO] validating %s", info.Name())
			if err = validateFile(path); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func validateFile(path string) error {
	filename := filepath.Base(path)
	//check that the file has an .xml extension
	if !xmlPtn.MatchString(filename) {
		log.Printf("[WARNING] file %s does not end in .xml skipping")
		return nil
	}
	return nil
}
