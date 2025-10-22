package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"wclogs-cli/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "🔧 Set up your Warcraft Logs API credentials",
	Long: color.HiCyanString(`
🔧 WARCRAFT LOGS CONFIG SETUP

This command helps you set up your Warcraft Logs API credentials interactively.
It will create a ~/.wclogs.yaml file with your client ID and secret.

To get your API credentials:
1. Go to https://www.warcraftlogs.com/api/clients
2. Create a new client (or use existing one)
3. Copy your Client ID and Client Secret

Your credentials will be stored in ~/.wclogs.yaml
`) + "\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigSetup()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfigSetup() error {
	reader := bufio.NewReader(os.Stdin)

	color.HiCyan("🔧 Warcraft Logs API Setup")
	color.HiCyan("========================\n")

	// Check if config already exists
	exists, err := config.ConfigExists()
	if err != nil {
		return fmt.Errorf("error checking config: %w", err)
	}

	if exists {
		fmt.Print("⚠️  Config file already exists. Overwrite? (y/N): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			color.HiGreen("✅ Config setup cancelled")
			return nil
		}
	}

	color.HiYellow("📋 Get your API credentials from:")
	color.HiYellow("   https://www.warcraftlogs.com/api/clients")
	fmt.Println()

	// Get Client ID
	fmt.Print("🔑 Enter your Client ID: ")
	clientID, _ := reader.ReadString('\n')
	clientID = strings.TrimSpace(clientID)

	if clientID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	// Get Client Secret
	fmt.Print("🔒 Enter your Client Secret: ")
	clientSecret, _ := reader.ReadString('\n')
	clientSecret = strings.TrimSpace(clientSecret)

	if clientSecret == "" {
		return fmt.Errorf("client secret cannot be empty")
	}

	// Create and save config
	cfg := &config.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Success message
	configPath, _ := config.GetConfigPath()
	color.HiGreen("✅ Configuration saved successfully!")
	color.HiGreen("📁 Config file: %s", configPath)
	fmt.Println()
	color.HiCyan("🚀 You can now use: wclogs damage <report> <fight>")

	return nil
}
