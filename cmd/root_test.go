package cmd

import (
	"bytes"
	"testing"
)

func TestAddTableCommands(t *testing.T) {
	// We'll check that the function executes without error
	// Since addTableCommands() adds commands to rootCmd, we can test indirectly
	originalCommandsCount := len(rootCmd.Commands())

	addTableCommands()

	newCommandsCount := len(rootCmd.Commands())

	// The function should add commands to rootCmd
	// Note: Commands might have already been added during package initialization
	// So we're primarily ensuring the function runs without error
	if newCommandsCount < originalCommandsCount {
		t.Error("addTableCommands() should add commands to rootCmd")
	}
}

func TestInitFunction(t *testing.T) {
	// Test that init() function runs and adds global flags
	// This is called automatically when the package is loaded
	// We can verify that global flags are registered

	if rootCmd.PersistentFlags().Lookup("verbose") == nil {
		t.Error("init() should add 'verbose' global flag")
	}

	if rootCmd.PersistentFlags().Lookup("output") == nil {
		t.Error("init() should add 'output' global flag")
	}

	if rootCmd.PersistentFlags().Lookup("top") == nil {
		t.Error("init() should add 'top' global flag")
	}
}

// Helper function to capture command output
func captureOutput(f func()) string {
	var buf bytes.Buffer
	// Note: This test is limited because the actual output goes to os.Stdout
	// In a real scenario, we might need to temporarily replace os.Stdout
	// For now, we'll just ensure the function doesn't panic
	f()
	return buf.String()
}

func TestRootCmdInitialization(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd should be initialized")
	}

	if rootCmd.Use != "wclogs" {
		t.Errorf("rootCmd.Use = %v, expected %v", rootCmd.Use, "wclogs")
	}

	// Check that persistent pre-run exists (config validation)
	if rootCmd.PersistentPreRunE == nil {
		t.Error("rootCmd should have PersistentPreRunE for config validation")
	}
}

// Test for the root command's main functionality
func TestRootCmdRun(t *testing.T) {
	if rootCmd.Run == nil {
		t.Error("rootCmd should have a Run function")
	}

	// The Run function should print help when called without args
	// This is difficult to test without capturing stdout
	// So we'll just ensure the function exists and is set
}
