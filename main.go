// Package main provides a minimal ls utility implementation.
package main

import (
	"fmt"
	"os"
	"sort"
)

// ls reads the directory and returns a list of non-hidden file names.
// If dir is empty, it defaults to ".".
func ls(dir string) ([]string, error) {
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, entry := range entries {
		name := entry.Name()
		if len(name) > 0 && name[0] != '.' {
			fileNames = append(fileNames, name)
		}
	}

	return fileNames, nil
}

// format sorts the file names alphabetically.
// C locale: case-sensitive ASCII order (e.g., LICENSE, Makefile, README.md, go-starter, go.mod, main.go)
// This differs from UTF-8 locale-aware sorting which is case-insensitive (e.g., go-starter, go.mod, LICENSE, main.go, Makefile, README.md)
func format(fileNames []string) []string {
	sorted := make([]string, len(fileNames))
	copy(sorted, fileNames)
	sort.Strings(sorted)
	return sorted
}

// printFiles prints each file name on a separate line.
func printFiles(fileNames []string) {
	for _, name := range fileNames {
		fmt.Println(name)
	}
}

// parseArgs parses command-line arguments and returns the long format flag and directory.
func parseArgs(args []string) (longFormat bool, dir string) {
	longFormat = false
	dir = "."

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "-l" {
			longFormat = true
		} else if arg[0] != '-' {
			// First non-flag argument is the directory
			dir = arg
			break
		}
	}

	return longFormat, dir
}

func main() {
	longFormat, dir := parseArgs(os.Args)
	_ = longFormat // Will be used in later steps for formatting

	fileNames, err := ls(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", dir, err)
		os.Exit(1)
	}

	formatted := format(fileNames)
	printFiles(formatted)
}
