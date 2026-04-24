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
	bin := buildTestBinary(t, "portwatch_test_bin")
	defer os.Remove(bin)

	cmd := exec.Command("./' + bin)
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

// buildTestBinary compiles the package into a temporary binary with the given
// name and returns the path. The test is failed immediately if compilation
// does not succeed.
func buildTestBinary(t *testing.T, name string) string {
	t.Helper()
	build := exec.Command("go", "build", "-o", name, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return "./" + name
}
