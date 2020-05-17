package main

import (
	"errors"
	"fmt"
	"maxsjolinder/go-peek/peinfo"
	"os"
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	filePath, err := getTargetFilePathArg()
	if err != nil {
		return err
	}

	peInfo, parseErr := peinfo.New(filePath)
	if parseErr != nil {
		return parseErr
	}

	peInfo.Print()
	return nil
}

func getTargetFilePathArg() (string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "", errors.New("No file specified. Please provide path to file to analyze")
	}
	return args[0], nil
}
