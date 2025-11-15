// Package main provides a minimal ls utility implementation.
package main

import (
	"fmt"
	"os"
)

func ls() {
	// Read current directory, print first file name
	entries, err := os.ReadDir(".")
	if err != nil {
		return // For now, silently fail if directory can't be read
	}
	if len(entries) > 0 {
		fmt.Println(entries[0].Name())
	}
}

func main() {
	ls()
}
