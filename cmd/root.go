package cmd

import (
	"fmt"
	"os"
	"strconv"

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
  wclogs interrupts ABC123 5  # Show interrupt events

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

	// Add all table commands - no separate files needed!
	addTableCommands()
}

// createTableHandler creates a command handler for the specified table type
func createTableHandler(tableType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Parse arguments
		reportCode := args[0]
		fightID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("fight-id must be a number, got: %s", args[1])
		}

		// Get flag values (inherited from root)
		topN, _ := cmd.Flags().GetInt("top")
		verbose, _ := cmd.Flags().GetBool("verbose")
		outputPath, _ := cmd.Flags().GetString("output")
		noColor, _ := cmd.Flags().GetBool("no-color")
		playerName, _ := cmd.Flags().GetString("player")

		// Call the shared handler with player filtering support
		return executeTableCommand(tableType, reportCode, fightID, topN, noColor, verbose, outputPath, playerName)
	}
}

// addTableCommands defines all table-based commands in one place
func addTableCommands() {
	// Damage command - WITH --player FLAG
	var damageCmd = &cobra.Command{
		Use:   "damage [report-code] [fight-id]",
		Short: "🗡️  Show damage table for a fight",
		Long: color.HiYellowString(`
🗡️  DAMAGE TABLE COMMAND

Display damage done by all players in a specific fight.

Examples:
  wclogs damage ABC123XYZ 5           # Show damage for fight 5
  wclogs damage ABC123XYZ 5 --top 10  # Show top 10 players only
  wclogs damage ABC123XYZ 5 --player "Pmpm"  # Show only specific player
  wclogs damage ABC123XYZ 5 --output damage.csv # Save to file
`) + "\n",
		Args: cobra.ExactArgs(2),
		RunE: createTableHandler("damage"),
	}
	damageCmd.Flags().BoolP("no-color", "n", false, "Disable color output")
	damageCmd.Flags().StringP("player", "p", "", "Filter by specific player name")
	rootCmd.AddCommand(damageCmd)

	// Healing command - NOW WITH --player FLAG
	var healingCmd = &cobra.Command{
		Use:   "healing [report-code] [fight-id]",
		Short: "💚 Show healing table for a fight",
		Long: color.HiGreenString(`
💚 HEALING TABLE COMMAND

Display healing done by all players in a specific fight.

Examples:
  wclogs healing ABC123XYZ 5           # Show healing for fight 5
  wclogs healing ABC123XYZ 5 --top 5   # Show top 5 healers only
  wclogs healing ABC123XYZ 5 --player "Sketch" # Show only specific player
  wclogs healing ABC123XYZ 5 --output healers.csv # Save to file
`) + "\n",
		Args: cobra.ExactArgs(2),
		RunE: createTableHandler("healing"),
	}
	healingCmd.Flags().BoolP("no-color", "n", false, "Disable color output")
	healingCmd.Flags().StringP("player", "p", "", "Filter by specific player name")
	rootCmd.AddCommand(healingCmd)

	// Deaths command - NOW WITH --player FLAG
	var deathsCmd = &cobra.Command{
		Use:   "deaths [report-code] [fight-id]",
		Short: "💀 Show death events for a fight",
		Long: color.HiRedString(`
💀 DEATHS TABLE COMMAND

Display death events for all players in a specific fight.

Examples:
  wclogs deaths ABC123XYZ 5            # Show deaths for fight 5
  wclogs deaths ABC123XYZ 5 --player "PlayerName" # Show specific player deaths
  wclogs deaths ABC123XYZ 5 --verbose  # Show detailed analysis
  wclogs deaths ABC123XYZ 5 --output deaths.json # Save to file
`) + "\n",
		Args: cobra.ExactArgs(2),
		RunE: createTableHandler("deaths"),
	}
	deathsCmd.Flags().BoolP("no-color", "n", false, "Disable color output")
	deathsCmd.Flags().StringP("player", "p", "", "Filter by specific player name")
	rootCmd.AddCommand(deathsCmd)

	// Interrupts command - NOW WITH --player FLAG
	var interruptsCmd = &cobra.Command{
		Use:   "interrupts [report-code] [fight-id]",
		Short: "🛑 Show interrupt events for a fight",
		Long: color.HiMagentaString(`
🛑 INTERRUPTS TABLE COMMAND

Display interrupt events performed by all players in a specific fight.

Examples:
  wclogs interrupts ABC123XYZ 5        # Show interrupts for fight 5
  wclogs interrupts ABC123XYZ 5 --top 10 # Show top interrupters
  wclogs interrupts ABC123XYZ 5 --player "PlayerName" # Show specific player
  wclogs interrupts ABC123XYZ 5 --output interrupts.csv # Save to file
`) + "\n",
		Args: cobra.ExactArgs(2),
		RunE: createTableHandler("interrupts"),
	}
	interruptsCmd.Flags().BoolP("no-color", "n", false, "Disable color output")
	interruptsCmd.Flags().StringP("player", "p", "", "Filter by specific player name")
	rootCmd.AddCommand(interruptsCmd)

	// Players command - NEW for Day 6 Afternoon
	var playersCmd = &cobra.Command{
		Use:   "players [report-code]",
		Short: "👥 Show all players in a report",
		Long: color.HiYellowString(`
👥 PLAYERS LIST COMMAND

Display all players who participated in a report, with their classes and servers.
This is useful for seeing who was in the raid and getting exact names for filtering.

Examples:
  wclogs players ABC123XYZ              # Show all players in report
  wclogs players ABC123XYZ --verbose    # Show detailed information
  wclogs players ABC123XYZ --output players.csv # Save to file

Use the exact names shown here with the --player flag:
  wclogs damage ABC123XYZ 5 --player "PlayerName"
`) + "\n",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportCode := args[0]
			verbose, _ := cmd.Flags().GetBool("verbose")
			outputPath, _ := cmd.Flags().GetString("output")
			return executePlayersCommand(reportCode, verbose, outputPath)
		},
	}
	rootCmd.AddCommand(playersCmd)
}
