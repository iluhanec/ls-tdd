package main

import (
	"bytes"
	"os"
	"testing"
)

func TestLsCommandExists(_ *testing.T) {
	// Test that ls command exists and can be invoked
	// This is the simplest possible test - just verify the function exists
	ls()
	// If we get here without panic, the command exists and can be invoked
}

func TestLsOutputsSomething(t *testing.T) {
	// Test that ls outputs something (even if just a placeholder)
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

	if output == "" {
		t.Error("ls() should output something, got empty string")
	}
}
