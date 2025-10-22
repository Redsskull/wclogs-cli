package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"wclogs-cli/config"
)

var rootCmd = &cobra.Command{
	Use:   "wclogs",
	Short: "🗡️  A CLI tool for Warcraft Logs analysis",
	Long: color.HiCyanString(`
🗡️  WARCRAFT LOGS CLI TOOL

A terminal-based tool for analyzing Warcraft Logs data using GraphQL.
Fast, scriptable access to combat log data without browser overhead.

Examples:
  wclogs damage ABC123 5      # Show damage table for fight 5
  wclogs healing ABC123 5     # Show healing table
  wclogs deaths ABC123 5      # Show death events
  wclogs players ABC123       # List all players in report

Get started by setting up your API credentials:
  wclogs config               # Interactive credential setup

This creates ~/.wclogs.yaml with your API keys.

For help with a specific command:
  wclogs help damage          # Help for damage command
`) + "\n",
	// Check for config before running any command that needs it
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config check for the config command itself and help
		if cmd.Name() == "config" || cmd.Name() == "help" {
			return nil
		}

		// Check if config exists
		exists, err := config.ConfigExists()
		if err != nil {
			return fmt.Errorf("error checking config: %w", err)
		}

		if !exists {
			color.HiRed("❌ No configuration found!")
			color.HiYellow("\n📋 Please set up your Warcraft Logs API credentials first:")
			color.HiWhite("   wclogs config")
			color.HiYellow("\nTo get API credentials:")
			color.HiYellow("   1. Go to https://www.warcraftlogs.com/api/clients")
			color.HiYellow("   2. Create a new client")
			color.HiYellow("   3. Run 'wclogs config' with your credentials")
			return fmt.Errorf("configuration required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		color.HiRed("❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags that work for all commands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Save output to file (format auto-detected from extension: .csv, .json)")
	rootCmd.PersistentFlags().IntP("top", "t", 0, "Show top N players (0 = all)")
}
