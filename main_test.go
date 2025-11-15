package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// createTestFiles creates test files in the specified directory.
func createTestFiles(t *testing.T, tmpDir string, fileNames []string) {
	t.Helper()
	for _, fileName := range fileNames {
		testFilePath := filepath.Clean(filepath.Join(tmpDir, fileName))
		testFile, err := os.Create(testFilePath)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", fileName, err)
		}
		if err := testFile.Close(); err != nil {
			t.Fatalf("failed to close test file %s: %v", fileName, err)
		}
	}
}

// runLsInTempDir changes to the workDir, captures stdout, runs ls(lsDir), and returns the output.
// If lsDir is empty, it defaults to "." (current directory).
// It restores the original directory and stdout before returning.
func runLsInTempDir(t *testing.T, workDir string, lsDir string) string {
	t.Helper()

	// Default lsDir to "." if not specified
	if lsDir == "" {
		lsDir = "."
	}

	// Save current working directory and change to workDir
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("failed to change to work directory: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	ls(lsDir)

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close write pipe: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	return buf.String()
}

func TestLsListsAtLeastOneFile(t *testing.T) {
	// Test that ls lists at least one file from current directory
	tmpDir := t.TempDir()
	testFileName := "testfile.txt"
	createTestFiles(t, tmpDir, []string{testFileName})

	output := runLsInTempDir(t, tmpDir, "")

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
	tmpDir := t.TempDir()

	// Create non-hidden files
	nonHiddenFiles := []string{"file1.txt", "file2.txt", "file3.go"}
	createTestFiles(t, tmpDir, nonHiddenFiles)

	// Create hidden files (should NOT be listed)
	hiddenFiles := []string{".hidden1", ".hidden2"}
	createTestFiles(t, tmpDir, hiddenFiles)

	output := runLsInTempDir(t, tmpDir, "")

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
	tmpDir := t.TempDir()

	// Create files in intentionally non-alphabetical order
	files := []string{"zebra.txt", "apple.txt", "banana.txt", "dog.txt"}
	createTestFiles(t, tmpDir, files)

	output := runLsInTempDir(t, tmpDir, "")

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

func TestLsAcceptsDirectoryArgument(t *testing.T) {
	// Test that ls <directory> lists files in specified directory
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create files in the parent directory (current directory during test)
	parentFiles := []string{"parent_file1.txt", "parent_file2.txt"}
	createTestFiles(t, tmpDir, parentFiles)

	// Create a subdirectory with its own files
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0750); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}
	subDirFiles := []string{"subfile1.txt", "subfile2.txt"}
	createTestFiles(t, subDir, subDirFiles)

	// Change to tmpDir and call ls(subDir)
	// This ensures that if ls() incorrectly reads the subDir instead of current directory,
	// it will find subDirFiles instead of parentFiles
	output := runLsInTempDir(t, tmpDir, subDir)

	// Verify output contains files from the specified subdirectory
	for _, fileName := range subDirFiles {
		if !bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should list file %q from directory %q, got output: %q", fileName, subDir, output)
		}
	}

	// Verify output does NOT contain files from current directory (tmpDir)
	// This proves we're reading the specified directory, not the current directory
	for _, fileName := range parentFiles {
		if bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should NOT list file %q from current directory %q, but it was found in output: %q. This indicates ls() is reading the wrong directory.", fileName, tmpDir, output)
		}
	}
}
