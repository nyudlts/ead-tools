package cmd

import (
	"encoding/xml"
	"fmt"
	"github.com/nyudlts/go-aspace"
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
		fmt.Println("validate-ead ==")
		if err := directoryExists(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(2)
		}

		fmt.Println("* validating ead files in %s", inputDir)
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
	fmt.Println("* validating files in %s", directory)
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
		log.Printf("[WARNING] file %s does not end in .xml skipping", filename)
		return nil
	}

	//get the bytes of a file
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[ERROR] could not read %s", filename)
		return nil
	}

	//check that ead is well-formed
	if err = xml.Unmarshal(fileBytes, new(interface{})); err != nil {
		log.Printf("[ERROR] %s is not well-formed", filename)
		return nil
	}

	//validate
	if err = aspace.ValidateEAD(fileBytes); err != nil {
		log.Printf("[ERROR] %s is not valid to EAD 2002 schema")
	}
	return nil
}
