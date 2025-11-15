package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLsListsAtLeastOneFile(t *testing.T) {
	// Test that ls lists at least one file from current directory
	// Create a temporary directory with a known test file
	tmpDir := t.TempDir()
	testFileName := "testfile.txt"
	testFilePath := filepath.Join(tmpDir, testFileName)
	// clean the file path to avoid any
	// potential security issues with path traversal
	testFilePath = filepath.Clean(testFilePath)

	// Create a test file in the temporary directory
	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := testFile.Close(); err != nil {
		t.Fatalf("failed to close test file: %v", err)
	}

	// Save current working directory and change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	ls()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close write pipe: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	output := buf.String()

	// Verify output contains the test file name
	if output == "" {
		t.Error("ls() should list at least one file, got empty output")
	}
	if output == "placeholder\n" {
		t.Error("ls() should list actual files, not placeholder")
	}
	// Check that the output contains our test file name
	if !bytes.Contains([]byte(output), []byte(testFileName)) {
		t.Errorf("ls() should list the test file %q, got output: %q", testFileName, output)
	}
}
