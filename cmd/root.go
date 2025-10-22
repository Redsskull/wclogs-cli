package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().StringP("format", "f", "table", "Output format (table, json, csv)")
	rootCmd.PersistentFlags().IntP("top", "t", 0, "Show top N players (0 = all)")
}
