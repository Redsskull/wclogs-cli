package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name:     "valid config",
			config:   Config{ClientID: "test_id", ClientSecret: "test_secret"},
			expected: true,
		},
		{
			name:     "empty client id",
			config:   Config{ClientID: "", ClientSecret: "test_secret"},
			expected: false,
		},
		{
			name:     "empty client secret",
			config:   Config{ClientID: "test_id", ClientSecret: ""},
			expected: false,
		},
		{
			name:     "both empty",
			config:   Config{ClientID: "", ClientSecret: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsValid()
			if result != tt.expected {
				t.Errorf("Config.IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Errorf("GetConfigPath() error = %v", err)
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Could not get home directory: %v", err)
	}

	expected := filepath.Join(home, ".wclogs.yaml")
	if path != expected {
		t.Errorf("GetConfigPath() = %v, expected %v", path, expected)
	}
}

func TestConfigExists(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempDir)

	// Test with non-existent config
	exists, err := ConfigExists()
	if err != nil {
		t.Errorf("ConfigExists() error = %v", err)
	}
	if exists {
		t.Error("ConfigExists() returned true for non-existent config")
	}

	// Create a config file
	configPath := filepath.Join(tempDir, ".wclogs.yaml")
	err = os.WriteFile(configPath, []byte("client_id: test"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test with existing config
	exists, err = ConfigExists()
	if err != nil {
		t.Errorf("ConfigExists() error = %v", err)
	}
	if !exists {
		t.Error("ConfigExists() returned false for existing config")
	}

	// Restore original home directory
	t.Setenv("HOME", originalHome)
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempDir)

	config := &Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
	}

	// Save config
	err := SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Load config
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if loadedConfig.ClientID != config.ClientID {
		t.Errorf("Loaded ClientID = %v, expected = %v", loadedConfig.ClientID, config.ClientID)
	}
	if loadedConfig.ClientSecret != config.ClientSecret {
		t.Errorf("Loaded ClientSecret = %v, expected = %v", loadedConfig.ClientSecret, config.ClientSecret)
	}

	// Restore original home directory
	t.Setenv("HOME", originalHome)
}

func TestLoadNonExistentConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempDir)

	// Test with non-existent config file
	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should return an error for non-existent config")
	}

	// Restore original home directory
	t.Setenv("HOME", originalHome)
}

func TestInvalidConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempDir)

	// Create an invalid config file (empty)
	configPath := filepath.Join(tempDir, ".wclogs.yaml")
	err := os.WriteFile(configPath, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Try to load the invalid config
	_, err = LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should return an error for invalid config")
	}

	// Restore original home directory
	t.Setenv("HOME", originalHome)
}