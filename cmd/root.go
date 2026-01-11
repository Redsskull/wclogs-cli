package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/models"
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
  wclogs damage ABC123 last   # Show damage for last fight
  wclogs healing ABC123 5     # Show healing table
  wclogs deaths ABC123 5      # Show death analysis

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
		fightIDStr := args[1]

		// Get verbose flag early for potential fight resolution logging
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Resolve fight ID (handles both numbers and "last" keyword)
		fightID, err := resolveFightID(reportCode, fightIDStr, verbose)
		if err != nil {
			return fmt.Errorf("failed to resolve fight ID '%s': %w", fightIDStr, err)
		}

		// Get flag values (inherited from root)
		topN, _ := cmd.Flags().GetInt("top")
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
		Use:   "damage [report-code] [fight-id|last]",
		Short: "🗡️  Show damage table for a fight",
		Long: color.HiYellowString(`
🗡️  DAMAGE TABLE COMMAND

Display damage done by all players in a specific fight.
Fight ID can be a number or "last" for the most recent fight.

Examples:
  wclogs damage ABC123XYZ 5           # Show damage for fight 5
  wclogs damage ABC123XYZ last        # Show damage for last fight
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
		Use:   "healing [report-code] [fight-id|last]",
		Short: "💚 Show healing table for a fight",
		Long: color.HiGreenString(`
💚 HEALING TABLE COMMAND

Display healing done by all players in a specific fight.
Fight ID can be a number or "last" for the most recent fight.

Examples:
  wclogs healing ABC123XYZ 5           # Show healing for fight 5
  wclogs healing ABC123XYZ last        # Show healing for last fight
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

	// Deaths Analysis command - Uses Events API for death analysis
	var deathsCmd = &cobra.Command{
		Use:   "deaths [report-code] [fight-id|last]",
		Short: "💀 Death analysis with summary and detailed modes",
		Long: color.HiRedString(`
💀 DEATH ANALYSIS

Two modes available:
• SUMMARY MODE (default): Concise overview of all deaths with timeline
• DETAILED MODE: In-depth analysis with healing/defensive data
Fight ID can be a number or "last" for the most recent fight.

Examples:
  wclogs deaths Hw9TZc2WyrVKJLCa 99                    # Summary of all deaths
  wclogs deaths Hw9TZc2WyrVKJLCa last                  # Summary for last fight
  wclogs deaths Hw9TZc2WyrVKJLCa 99 --player "Jusdis"  # Detailed analysis for specific player
  wclogs deaths Hw9TZc2WyrVKJLCa 99 --verbose          # Verbose summary mode
`) + "\n",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			playerName, _ := cmd.Flags().GetString("player")
			return ExecuteDeathAnalysis(args[0], args[1], playerName, verbose)
		},
	}
	deathsCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	deathsCmd.Flags().StringP("player", "p", "", "Filter to specific player")
	rootCmd.AddCommand(deathsCmd)

	// Interrupt Analysis command - Uses Events API for interrupt analysis
	var interruptCmd = &cobra.Command{
		Use:   "interrupts [report-code] [fight-id|last]",
		Short: "🎛️  Interrupt analysis with detailed breakdown",
		Long: color.HiBlueString(`
🎛️  INTERRUPT ANALYSIS

Analyze interrupts performed and casts that were stopped during a fight.
Fight ID can be a number or "last" for the most recent fight.

Examples:
  wclogs interrupts Hw9TZc2WyrVKJLCa 99                    # Summary of all interrupts
  wclogs interrupts Hw9TZc2WyrVKJLCa last                  # Summary for last fight
  wclogs interrupts Hw9TZc2WyrVKJLCa 99 --player "PlayerName"  # Detailed analysis for specific player
  wclogs interrupts Hw9TZc2WyrVKJLCa 99 --verbose          # Verbose interrupt analysis
`) + "\n",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			playerName, _ := cmd.Flags().GetString("player")
			return ExecuteInterruptAnalysis(args[0], args[1], playerName, verbose)
		},
	}
	interruptCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	interruptCmd.Flags().StringP("player", "p", "", "Filter to specific player")
	rootCmd.AddCommand(interruptCmd)

	// Players command - List all players in a report
	var playersCmd = &cobra.Command{
		Use:   "players [report-code] [fight-id]",
		Short: "👥 List all players in a report or specific fight",
		Long: color.HiMagentaString(`
👥 PLAYERS COMMAND

List all players in a report with their classes, servers, and item levels.
Supports filtering by class, role, or player name for easy searching.
Optionally filter to players who participated in a specific fight.

Examples:
  wclogs players ABC123XYZ                      # List all players in report
  wclogs players ABC123XYZ 5                    # List players in fight 5
  wclogs players ABC123XYZ last                 # List players in last fight
  wclogs players ABC123XYZ --class "Paladin"    # Filter by class
  wclogs players ABC123XYZ --role "Tank"        # Filter by role
  wclogs players ABC123XYZ --search "Pmpm"      # Search by player name
  wclogs players ABC123XYZ --output players.csv # Export to file
`) + "\n",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			classFilter, _ := cmd.Flags().GetString("class")
			roleFilter, _ := cmd.Flags().GetString("role")
			searchFilter, _ := cmd.Flags().GetString("search")
			outputPath, _ := cmd.Flags().GetString("output")
			topN, _ := cmd.Flags().GetInt("top")
			noColor, _ := cmd.Flags().GetBool("no-color")

			// Handle optional fight ID parameter
			var fightIDStr string
			if len(args) > 1 {
				fightIDStr = args[1]
			}

			debug, _ := cmd.Flags().GetBool("debug")
			return ExecutePlayersCommand(args[0], fightIDStr, classFilter, roleFilter, searchFilter, outputPath, topN, noColor, verbose, debug)
		},
	}
	playersCmd.Flags().BoolP("no-color", "n", false, "Disable color output")
	playersCmd.Flags().StringP("class", "c", "", "Filter by class (e.g., Paladin, Warrior)")
	playersCmd.Flags().StringP("role", "r", "", "Filter by role (Tank, Healer, DPS)")
	playersCmd.Flags().StringP("search", "s", "", "Search by player name")
	playersCmd.Flags().BoolP("debug", "d", false, "Show detailed spec icon debugging information")
	rootCmd.AddCommand(playersCmd)
}

// resolveFightID resolves a fight ID string to an actual numeric fight ID
// Supports both numeric IDs and the "last" keyword for the last fight in the report
func resolveFightID(reportCode string, fightIDStr string, verbose bool) (int, error) {
	// If it's already a number, parse and return it
	if fightID, err := strconv.Atoi(fightIDStr); err == nil {
		return fightID, nil
	}

	// Handle the "last" keyword
	if fightIDStr == "last" {
		if verbose {
			color.HiBlue("🔍 Resolving 'last' fight ID for report %s...", reportCode)
		}

		// Setup API client for fight resolution
		cfg, err := config.LoadConfig()
		if err != nil {
			return 0, fmt.Errorf("failed to load config: %w", err)
		}

		authClient := auth.NewClient(cfg.ClientID, cfg.ClientSecret)
		apiClient := api.NewClient(authClient)

		// Query for fight information
		fightRequest := api.NewFightInfoRequest(reportCode)
		response, err := apiClient.Query(fightRequest.Query, fightRequest.Variables)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch fight info: %w", err)
		}

		// Parse fights from response
		if response.Data == nil || response.Data.ReportData == nil || response.Data.ReportData.Report == nil {
			return 0, fmt.Errorf("no fight data found in report")
		}

		fights := response.Data.ReportData.Report.Fights
		if len(fights) == 0 {
			return 0, fmt.Errorf("no fights found in report")
		}

		// Find the last meaningful fight (not just chronologically last)
		lastMeaningfulFight := findLastMeaningfulFight(fights, verbose)
		if lastMeaningfulFight == nil {
			// Fallback to actual last fight if no meaningful fight found
			lastMeaningfulFight = &fights[len(fights)-1]
			if verbose {
				color.HiYellow("⚠️  No meaningful fights found, using chronologically last fight")
			}
		}

		if verbose {
			color.HiGreen("✅ Resolved 'last' to fight #%d: %s", lastMeaningfulFight.ID, lastMeaningfulFight.Name)
		}

		return lastMeaningfulFight.ID, nil
	}

	// Invalid fight ID format
	return 0, fmt.Errorf("fight-id must be a number or 'last', got: %s", fightIDStr)
}

// findLastMeaningfulFight finds the last fight that represents a meaningful boss encounter
// Uses official WCL logic: fights with encounterID = 0 are considered trash fights
func findLastMeaningfulFight(fights []models.Fight, verbose bool) *models.Fight {
	if len(fights) == 0 {
		return nil
	}

	var lastMeaningfulFight *models.Fight

	// Go through fights in reverse order to find the last meaningful one
	for i := len(fights) - 1; i >= 0; i-- {
		fight := &fights[i]

		// Official WCL logic: encounterID = 0 means trash fight
		if fight.EncounterID == 0 {
			if verbose {
				color.HiBlack("⏭️  Skipping fight #%d (%s) - trash fight (encounterID = 0)",
					fight.ID, fight.Name)
			}
			continue
		}

		// This is a meaningful encounter (boss fight)
		lastMeaningfulFight = fight
		if verbose {
			killStatus := "WIPE"
			if fight.Kill {
				killStatus = "KILL"
			}
			duration := (fight.EndTime - fight.StartTime) / 1000 // Convert to seconds
			color.HiGreen("✅ Found last meaningful fight #%d: %s (%s, %.1fs, encounterID: %d)",
				fight.ID, fight.Name, killStatus, duration, fight.EncounterID)
		}
		break
	}

	return lastMeaningfulFight
}
