package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// IsValid checks if the config has the required fields
func (c *Config) IsValid() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

// GetConfigPath returns the path to the config file (~/.wclogs.yaml)
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	return filepath.Join(home, ".wclogs.yaml"), nil
}

// LoadConfig loads configuration from ~/.wclogs.yaml
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s\n\nPlease run 'wclogs config' to set up your credentials", configPath)
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	// Validate
	if !config.IsValid() {
		return nil, fmt.Errorf("invalid config file: client_id and client_secret are required")
	}

	return &config, nil
}

// SaveConfig saves the configuration to ~/.wclogs.yaml
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create YAML data
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot create YAML: %w", err)
	}

	// Write to file with appropriate permissions (user read/write only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}

	return nil
}

// ConfigExists checks if the config file exists
func ConfigExists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	return !os.IsNotExist(err), nil
}
