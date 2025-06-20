package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCommandExecution tests that commands run and return correct output
func TestCommandExecution(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "echo command",
			args:           []string{"run", "echo", "Hello World"},
			expectedOutput: "Hello World\n",
			expectedError:  false,
		},
		{
			name:           "pwd command",
			args:           []string{"run", "pwd"},
			expectedOutput: "gocker\n",
			expectedError:  false,
		},
		{
			name:          "invalid command",
			args:          []string{"run", "noncommand"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", append([]string{"run", "main.go"}, tt.args...)...)
			output, err := cmd.CombinedOutput()

			if tt.expectedError {
				if err == nil {
					t.Errorf("No expected error.")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v.", err)
				return
			}

			// Check only the suffix for the pwd command (also enough for echo)
			if !strings.HasSuffix(string(output), tt.expectedOutput) {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, string(output))
			}
		})
	}
}

// TestExitCodes verifies that exit codes are properly propagated
func TestExitCodes(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedCode int
	}{
		{
			name:         "successful command",
			args:         []string{"run", "true"},
			expectedCode: 0,
		},
		{
			name:         "failing command",
			args:         []string{"run", "false"},
			expectedCode: 1,
		},
		{
			name:         "ls nonexistent file",
			args:         []string{"run", "ls", "nonexistent_file"},
			expectedCode: 2, // "No such file or directory"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Needs a compiled binary for proper error propagation
			cmd := exec.Command("./gocker", tt.args...)
			err := cmd.Run()

			var exitCode int
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					t.Fatalf("Unexpected error type: %v", err)
				}
			}

			if exitCode != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectedCode, exitCode)
			}
		})
	}
}

// TestNamespaceIsolation tests that UTS namespace is properly isolated
func TestNamespaceIsolation(t *testing.T) {
	// Skip test if not running as root (required for namespaces)
	if os.Geteuid() != 0 {
		t.Skip("Skipping namespace test (requires root privileges)")
	}

	// Get original hostname
	originalHostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Failed to get original hostname: %v", err)
	}

	// Test that container has different hostname
	cmd := exec.Command("go", "run", "main.go", "run", "hostname")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run container: %v", err)
	}

	containerHostname := strings.TrimSpace(string(output))
	if containerHostname != "container" {
		t.Errorf("Expected container hostname 'container', got %q", containerHostname)
	}

	// Verify host hostname unchanged
	currentHostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("Failed to get current hostname: %v", err)
	}

	if currentHostname != originalHostname {
		t.Errorf("Host hostname changed from %q to %q", originalHostname, currentHostname)
	}
}
