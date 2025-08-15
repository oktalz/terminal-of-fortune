package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func readFiles() {
	var err error
	var files []string
	if len(os.Args) == 2 {
		files = []string{os.Args[1]}
	}

	if len(files) == 0 {
		files, err = filepath.Glob("*.txt")
		if err != nil {
			fmt.Println("Error finding .txt files:", err)
			os.Exit(1)
		}
	}

	var userFile string
	switch len(files) {
	case 0:
		fmt.Println("No .txt files found in the current directory.")
		os.Exit(1)
	case 1:
		userFile = files[0]
	default:
		userFile, err = runChooser(files)
		if err != nil {
			fmt.Println("Error running chooser:", err)
			os.Exit(1)
		}
		if userFile == "" {
			os.Exit(0)
		}
	}

	if err := readUsers(userFile); err != nil {
		fmt.Println("Error reading users:", err)
		os.Exit(1)
	}
}
