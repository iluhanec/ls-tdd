// Package main provides a minimal ls utility implementation.
package main

import (
	"fmt"
	"os"
	"sort"
)

func ls() {
	// Read current directory, print all non-hidden files
	entries, err := os.ReadDir(".")
	if err != nil {
		return // For now, silently fail if directory can't be read
	}

	// Collect non-hidden file names
	var fileNames []string
	for _, entry := range entries {
		// Skip hidden files (files starting with .)
		name := entry.Name()
		if len(name) > 0 && name[0] != '.' {
			fileNames = append(fileNames, name)
		}
	}

	// Sort file names alphabetically (C locale: case-sensitive ASCII order)
	sort.Strings(fileNames)

	// Print sorted file names
	for _, name := range fileNames {
		fmt.Println(name)
	}
}

func main() {
	ls()
}
