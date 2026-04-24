package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestMainBinaryBuilds verifies the package compiles without errors.
func TestMainBinaryBuilds(t *testing.T) {
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	// Clean up the produced binary if any.
	os.Remove("portwatch")
}

// TestMainExitsOnSIGTERM starts the binary and sends SIGTERM, expecting a
// clean exit within a reasonable timeout.
func TestMainExitsOnSIGTERM(t *testing.T) {
	// Build first.
	build := exec.Command("go", "build", "-o", "portwatch_test_bin", ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	defer os.Remove("portwatch_test_bin")

	cmd := exec.Command("./portwatch_test_bin")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start binary: %v", err)
	}

	// Give it a moment to initialise.
	time.Sleep(200 * time.Millisecond)

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("failed to send interrupt: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-done:
		// Process exited — success.
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
		t.Fatal("binary did not exit within timeout after SIGTERM")
	}
}
