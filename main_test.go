package main

import "testing"

func TestLsCommandExists(_ *testing.T) {
	// Test that ls command exists and can be invoked
	// This is the simplest possible test - just verify the function exists
	ls()
	// If we get here without panic, the command exists and can be invoked
}
