package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

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

// runLsInTempDir changes to workDir, captures stdout, runs ls(lsDir), and returns the output.
// If lsDir is empty, it defaults to ".".
func runLsInTempDir(t *testing.T, workDir string, lsDir string) string {
	t.Helper()

	if lsDir == "" {
		lsDir = "."
	}

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

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	_ = ls(lsDir)

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
	tmpDir := t.TempDir()
	testFileName := "testfile.txt"
	createTestFiles(t, tmpDir, []string{testFileName})

	output := runLsInTempDir(t, tmpDir, "")

	if output == "" {
		t.Error("ls() should list at least one file, got empty output")
	}
	if output == "placeholder\n" {
		t.Error("ls() should list actual files, not placeholder")
	}
	if !bytes.Contains([]byte(output), []byte(testFileName)) {
		t.Errorf("ls() should list the test file %q, got output: %q", testFileName, output)
	}
}

func TestLsListsAllNonHiddenFiles(t *testing.T) {
	tmpDir := t.TempDir()

	nonHiddenFiles := []string{"file1.txt", "file2.txt", "file3.go"}
	createTestFiles(t, tmpDir, nonHiddenFiles)

	hiddenFiles := []string{".hidden1", ".hidden2"}
	createTestFiles(t, tmpDir, hiddenFiles)

	output := runLsInTempDir(t, tmpDir, "")

	for _, fileName := range nonHiddenFiles {
		if !bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should list non-hidden file %q, got output: %q", fileName, output)
		}
	}

	for _, fileName := range hiddenFiles {
		if bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should NOT list hidden file %q, but it was found in output: %q", fileName, output)
		}
	}

	if output == "" {
		t.Error("ls() should list at least one file, got empty output")
	}
}

func TestLsSortsOutputAlphabetically(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"zebra.txt", "apple.txt", "banana.txt", "dog.txt"}
	createTestFiles(t, tmpDir, files)

	output := runLsInTempDir(t, tmpDir, "")

	lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))
	expectedOrder := []string{"apple.txt", "banana.txt", "dog.txt", "zebra.txt"}

	if len(lines) != len(expectedOrder) {
		t.Errorf("expected %d files, got %d. Output: %q", len(expectedOrder), len(lines), output)
		return
	}

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
	tmpDir := t.TempDir()

	parentFiles := []string{"parent_file1.txt", "parent_file2.txt"}
	createTestFiles(t, tmpDir, parentFiles)

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0750); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}
	subDirFiles := []string{"subfile1.txt", "subfile2.txt"}
	createTestFiles(t, subDir, subDirFiles)

	output := runLsInTempDir(t, tmpDir, subDir)

	for _, fileName := range subDirFiles {
		if !bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should list file %q from directory %q, got output: %q", fileName, subDir, output)
		}
	}

	for _, fileName := range parentFiles {
		if bytes.Contains([]byte(output), []byte(fileName)) {
			t.Errorf("ls() should NOT list file %q from current directory %q, but it was found in output: %q. This indicates ls() is reading the wrong directory.", fileName, tmpDir, output)
		}
	}
}

// runLsWithErrorCapture captures both stdout and stderr, runs ls(lsDir), and returns both outputs.
func runLsWithErrorCapture(t *testing.T, workDir string, lsDir string) (stdout string, stderr string) {
	t.Helper()

	if lsDir == "" {
		lsDir = "."
	}

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

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	os.Stdout = stdoutW

	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}
	os.Stderr = stderrW

	_ = ls(lsDir)

	if err := stdoutW.Close(); err != nil {
		t.Fatalf("failed to close stdout write pipe: %v", err)
	}
	if err := stderrW.Close(); err != nil {
		t.Fatalf("failed to close stderr write pipe: %v", err)
	}

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var stdoutBuf bytes.Buffer
	if _, err := stdoutBuf.ReadFrom(stdoutR); err != nil {
		t.Fatalf("failed to read from stdout pipe: %v", err)
	}

	var stderrBuf bytes.Buffer
	if _, err := stderrBuf.ReadFrom(stderrR); err != nil {
		t.Fatalf("failed to read from stderr pipe: %v", err)
	}

	return stdoutBuf.String(), stderrBuf.String()
}

func TestLsHandlesInvalidDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "nonexistent_directory")

	stdout, stderr := runLsWithErrorCapture(t, tmpDir, invalidPath)

	if stderr == "" {
		t.Error("ls() should print an error message to stderr for invalid directory, got empty stderr")
	}

	if !bytes.Contains([]byte(stderr), []byte("nonexistent_directory")) {
		t.Errorf("ls() error message should mention the invalid path, got stderr: %q", stderr)
	}

	if stdout != "" {
		t.Errorf("ls() should not print anything to stdout for invalid directory, got: %q", stdout)
	}
}
