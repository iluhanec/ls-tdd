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

func TestLsListsAllNonHiddenFiles(t *testing.T) {
	// Test that ls lists all non-hidden files in current directory
	// Create a temporary directory with multiple files (some hidden, some not)
	tmpDir := t.TempDir()

	// Create non-hidden files
	nonHiddenFiles := []string{"file1.txt", "file2.txt", "file3.go"}
	for _, fileName := range nonHiddenFiles {
		testFilePath := filepath.Clean(filepath.Join(tmpDir, fileName))
		testFile, err := os.Create(testFilePath)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", fileName, err)
		}
		if err := testFile.Close(); err != nil {
			t.Fatalf("failed to close test file %s: %v", fileName, err)
		}
	}

	// Create hidden files (should NOT be listed)
	hiddenFiles := []string{".hidden1", ".hidden2"}
	for _, fileName := range hiddenFiles {
		testFilePath := filepath.Clean(filepath.Join(tmpDir, fileName))
		testFile, err := os.Create(testFilePath)
		if err != nil {
			t.Fatalf("failed to create hidden file %s: %v", fileName, err)
		}
		if err := testFile.Close(); err != nil {
			t.Fatalf("failed to close hidden file %s: %v", fileName, err)
		}
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

	// Verify all non-hidden files are listed
	for _, fileName := range nonHiddenFiles {
		if !bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should list non-hidden file %q, got output: %q", fileName, output)
		}
	}

	// Verify hidden files are NOT listed
	for _, fileName := range hiddenFiles {
		if bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should NOT list hidden file %q, but it was found in output: %q", fileName, output)
		}
	}

	// Verify we have output (at least one file listed)
	if output == "" {
		t.Error("ls() should list at least one file, got empty output")
	}
}

func TestLsSortsOutputAlphabetically(t *testing.T) {
	// Test that files are listed in alphabetical order
	// Create a temporary directory with files in non-alphabetical order
	tmpDir := t.TempDir()

	// Create files in intentionally non-alphabetical order
	files := []string{"zebra.txt", "apple.txt", "banana.txt", "dog.txt"}
	for _, fileName := range files {
		testFilePath := filepath.Clean(filepath.Join(tmpDir, fileName))
		testFile, err := os.Create(testFilePath)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", fileName, err)
		}
		if err := testFile.Close(); err != nil {
			t.Fatalf("failed to close test file %s: %v", fileName, err)
		}
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

	// Split output into lines and remove trailing newline
	lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))

	// Expected alphabetical order
	expectedOrder := []string{"apple.txt", "banana.txt", "dog.txt", "zebra.txt"}

	// Verify we have the correct number of files
	if len(lines) != len(expectedOrder) {
		t.Errorf("expected %d files, got %d. Output: %q", len(expectedOrder), len(lines), output)
		return
	}

	// Verify files are in alphabetical order
	for i, expectedName := range expectedOrder {
		if i >= len(lines) {
			t.Errorf("not enough lines in output. Expected %d, got %d", len(expectedOrder), len(lines))
			break
		}
		actualName := string(lines[i])
		if actualName != expectedName {
			t.Errorf("file at position %d should be %q, got %q. Full output: %q", i, expectedName, actualName, output)
		}
	}
}
