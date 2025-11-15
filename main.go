// Package main provides a minimal ls utility implementation.
package main

import (
	"fmt"
	"os"
)

func ls() {
	// Read current directory, print all non-hidden files
	entries, err := os.ReadDir(".")
	if err != nil {
		return // For now, silently fail if directory can't be read
	}
	for _, entry := range entries {
		// Skip hidden files (files starting with .)
		name := entry.Name()
		if len(name) > 0 && name[0] != '.' {
			fmt.Println(name)
		}
	}
}

func main() {
	ls()
}
