// Package main provides a minimal ls utility implementation.
package main

import (
	"fmt"
	"os"
	"sort"
)

func ls(dir string) error {
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", dir, err)
		return err
	}

	var fileNames []string
	for _, entry := range entries {
		name := entry.Name()
		if len(name) > 0 && name[0] != '.' {
			fileNames = append(fileNames, name)
		}
	}

	// C locale: case-sensitive ASCII order (e.g., LICENSE, Makefile, README.md, go-starter, go.mod, main.go)
	// This differs from UTF-8 locale-aware sorting which is case-insensitive (e.g., go-starter, go.mod, LICENSE, main.go, Makefile, README.md)
	sort.Strings(fileNames)

	for _, name := range fileNames {
		fmt.Println(name)
	}
	return nil
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if err := ls(dir); err != nil {
		os.Exit(1)
	}
}
